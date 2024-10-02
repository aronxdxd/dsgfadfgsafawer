package global

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	payload := map[string]interface{}{"ok": false, "error": message}
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	payloadMap := map[string]interface{}{}
	if originalPayload, ok := payload.(map[string]interface{}); ok {
		for k, v := range originalPayload {
			payloadMap[k] = v
		}
	} else {
		payloadMap["data"] = payload
	}
	payloadMap["ok"] = true
	response, _ := json.Marshal(payloadMap)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
