// Start of Selection
package funcs

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"project/db/operations"
	"project/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		global.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var userInfo global.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		global.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if userInfo.TelegramID == 0 {
		global.RespondWithError(w, http.StatusBadRequest, "telegram_id is required")
		return
	}

	user, err := operations.Get(r.Context(), userInfo.TelegramID, bson.M{"user_id": userInfo.TelegramID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			global.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			log.Printf("Error finding user: %v", err)
			global.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	referrals, ok := user["referrals"].(primitive.A)
	if !ok {
		log.Printf("Error: referrals is not of type primitive.A")
		global.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	user["referrals"] = len(referrals)
	
	adjustedBalance := float64(user["balance"].(int32)) * 0.00449986 
	user["converted_balance"] = adjustedBalance

	global.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"ok":   true,
		"info": user,
	})
	log.Println("Successfully processed request in /get_info")
}