package modifies

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"

	"project/db/operations"
	"project/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleBuyToques(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		global.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var buyInfo global.BuyToquesInfo
	if err := json.NewDecoder(r.Body).Decode(&buyInfo); err != nil {
		global.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if buyInfo.TelegramID == 0 {
		global.RespondWithError(w, http.StatusBadRequest, "telegram_id is required")
		return
	}

	ctx := r.Context()
	user, err := operations.Get(ctx, buyInfo.TelegramID, bson.M{"user_id": buyInfo.TelegramID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			global.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			log.Printf("Error finding user: %v", err)
			global.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	fmt.Println(user)
	balance, ok := user["balance"].(int32)
	if !ok {
		global.RespondWithError(w, http.StatusInternalServerError, "Invalid balance data")
		return
	}
	modifies, ok := user["modifies"].(map[string]interface{})
	if !ok {
		global.RespondWithError(w, http.StatusInternalServerError, "Invalid modifies data")
		return
	}

	userLevel, ok := modifies["toques_lvl"].(int32)
	if !ok {
		global.RespondWithError(w, http.StatusInternalServerError, "Invalid toques level data")
		return
	}

	if userLevel >= int32(buyInfo.Level) || userLevel+1 != int32(buyInfo.Level) {
		global.RespondWithError(w, http.StatusBadRequest, "Invalid level request")
		return
	}

	var priceLevel int32

	switch buyInfo.Level {
	case 2:
		priceLevel = 1000
	default:
		priceLevel = int32(1500 * math.Pow(1.5, float64(buyInfo.Level-3)))
	}


	log.Printf("User %d is buying toques level %d. Current balance: %d, Price: %d", buyInfo.TelegramID, buyInfo.Level, balance, priceLevel)

	if balance < priceLevel {
		global.RespondWithError(w, http.StatusBadRequest, "Insufficient balance")
		return
	}

	if err := operations.IncrementUserField(ctx, buyInfo.TelegramID, "modifies.toques_lvl", 1); err != nil {
		log.Printf("Failed to update toques level for user %d: %v", buyInfo.TelegramID, err)
		global.RespondWithError(w, http.StatusInternalServerError, "Failed to update toques level")
		return
	}

	if err := operations.IncrementUserField(ctx, buyInfo.TelegramID, "balance", -priceLevel); err != nil {
		log.Printf("Failed to update balance for user %d: %v", buyInfo.TelegramID, err)
		global.RespondWithError(w, http.StatusInternalServerError, "Failed to update balance")
		return
	}

	global.RespondWithJson(w, http.StatusOK, map[string]bool{"ok": true})
}
