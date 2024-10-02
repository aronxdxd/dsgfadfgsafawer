package global

type UserInfo struct {
	TelegramID int `json:"telegram_id"`
}

type ClickInfo struct {
	TelegramID int `json:"telegram_id"`
	ClickCount int `json:"click_count"`
}

type BuyToquesInfo struct {
	TelegramID int `json:"telegram_id"`
	Level      int `json:"level"`
}