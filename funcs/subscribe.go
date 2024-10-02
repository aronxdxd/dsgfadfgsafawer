package funcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"project/db/operations"
	"project/global"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type telegramAPI struct {
	client  *http.Client
	token   string
	chatID  string
	baseURL string
}

func newTelegramAPI(token, chatID string) *telegramAPI {
	return &telegramAPI{
		client:  &http.Client{},
		token:   token,
		chatID:  chatID,
		baseURL: "https://api.telegram.org/bot",
	}
}

func (t *telegramAPI) CheckMembership(userID int) (bool, error) {
	url := fmt.Sprintf("%s%s/getChatMember", t.baseURL, t.token)
	payload := map[string]interface{}{
		"chat_id": t.chatID,
		"user_id": userID,
	}

	resp, err := t.makeRequest(url, payload)
	if err != nil {
		return false, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	var result struct {
		Ok     bool `json:"ok"`
		Result struct {
			Status string `json:"status"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, errors.Wrap(err, "failed to parse API response")
	}

	return result.Ok && (result.Result.Status == "member" || result.Result.Status == "administrator" || result.Result.Status == "creator"), nil
}

func (t *telegramAPI) makeRequest(url string, payload interface{}) (*http.Response, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal payload")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	return resp, nil
}

func CheckSubscription(w http.ResponseWriter, r *http.Request) {
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

	config, err := global.LoadConfig()
	if err != nil {
		global.RespondWithError(w, http.StatusInternalServerError, "Failed to load configuration")
		return
	}

	checker := newTelegramAPI(config.Token, config.ChatID)
	subscribed, err := checker.CheckMembership(userInfo.TelegramID)
	if err != nil {
		global.RespondWithError(w, http.StatusInternalServerError, "Failed to check subscription")
		return
	}

	user, err := operations.Get(context.Background(), userInfo.TelegramID, bson.M{"user_id": userInfo.TelegramID})
	if err != nil {
		global.RespondWithError(w, http.StatusInternalServerError, "Failed to get user info")
		return
	}

	fmt.Println(user)
	userModifies, ok := user["modifies"].(map[string]interface{})
	if !ok {
		global.RespondWithError(w, http.StatusInternalServerError, "Failed to get user modifies")
		return
	}


	if subscribed {
		if !userModifies["subscribe_bonus"].(bool) {
			userModifies["subscribe_bonus"] = true
			operations.UpdateUserField(context.Background(), userInfo.TelegramID, "modifies", userModifies)

			operations.IncrementUserField(context.Background(), userInfo.TelegramID, "energy_max", 1000)
		}
	}

	global.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"subscribed": subscribed,
	})
}
