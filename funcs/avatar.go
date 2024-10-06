package funcs

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func GetAvatar(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	telegramID := r.FormValue("telegram_id")
	if telegramID == "" {
		http.Error(w, "telegram_id is required", http.StatusBadRequest)
		return
	}

	_, err := strconv.ParseInt(telegramID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid telegram_id format", http.StatusBadRequest)
		return
	}
	avatarPath := "./avatar/" + telegramID + ".png"
	fmt.Println("Attempting to serve avatar from:", avatarPath)
	
	_, err = os.Stat(avatarPath)
	if os.IsNotExist(err) {
		http.Error(w, "Avatar not found", http.StatusNotFound)
		return
	}
	
	http.ServeFile(w, r, avatarPath)
}
