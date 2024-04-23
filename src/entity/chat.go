package entity

type Chat struct {
	Type    int    `json:"type"`
	Sender  int    `json:"sender"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    []*Chat `json:"data"`
}

type SendChatRequest struct {
	UserId     string `json:"userId"`
	LocationId string `json:"locationId"`
	Content    string `json:"content"`
}
