package funcs

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"project/db/operations"
	"project/global"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	errMethodNotAllowed = errors.New("only POST method is allowed")
	errInvalidPayload   = errors.New("invalid request payload")
	errTelegramIDRequired = errors.New("telegram_id is required")
	errClickCountRequired = errors.New("click_count is required")
	errGetUserInfo      = errors.New("failed to get user info")
	errInvalidEnergy    = errors.New("invalid energy value")
	errInsufficientEnergy = errors.New("not enough energy")
	errUpdateEnergy     = errors.New("failed to update user energy")
	errUpdateBalance    = errors.New("failed to update user balance")
	errInvalidToquesLvl = errors.New("invalid toques")
)

func HandleClicker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		global.RespondWithError(w, http.StatusMethodNotAllowed, errMethodNotAllowed.Error())
		return
	}

	var clickInfo global.ClickInfo
	if err := json.NewDecoder(r.Body).Decode(&clickInfo); err != nil {
		global.RespondWithError(w, http.StatusBadRequest, errInvalidPayload.Error())
		return
	}

	if err := validateClickInfo(clickInfo); err != nil {
		global.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedEnergy, err := processClick(r.Context(), clickInfo)
	if err != nil {
		global.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	global.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"now_energy": updatedEnergy,
	})
}

func validateClickInfo(clickInfo global.ClickInfo) error {
	switch {
	case clickInfo.TelegramID == 0:
		return errTelegramIDRequired
	case clickInfo.ClickCount == 0:
		return errClickCountRequired
	default:
		return nil
	}
}

func processClick(ctx context.Context, clickInfo global.ClickInfo) (int, error) {
	user, err := operations.Get(ctx, clickInfo.TelegramID, bson.M{"user_id": clickInfo.TelegramID})
	if err != nil {
		return 0, errGetUserInfo
	}

	userEnergy, ok := user["energy"].(int32)
	if !ok {
		return 0, errInvalidEnergy
	}

	if int64(userEnergy) < int64(clickInfo.ClickCount) {
		return 0, errInsufficientEnergy
	}

	modifies, ok := user["modifies"].(map[string]interface{})
	if !ok {
		return 0, errInvalidToquesLvl
	}

	toquesLvl, ok := modifies["toques_lvl"].(int32)
	if !ok {
		return 0, errInvalidToquesLvl
	}

	balance := clickInfo.ClickCount * int(toquesLvl)

	if err := operations.IncrementUserField(ctx, clickInfo.TelegramID, "energy", -clickInfo.ClickCount); err != nil {
		return 0, errUpdateEnergy
	}

	if err := operations.IncrementUserField(ctx, clickInfo.TelegramID, "balance", balance); err != nil {
		return 0, errUpdateBalance
	}

	updatedEnergy := int(userEnergy) - clickInfo.ClickCount
	log.Printf("Energy and balance updated for TelegramID %d. Old energy: %d, New energy: %d", clickInfo.TelegramID, userEnergy, updatedEnergy)
	return updatedEnergy, nil
}