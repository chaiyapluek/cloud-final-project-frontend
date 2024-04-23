package controller

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"

	"dev.chaiyapluek.cloud.final.frontend/src/service"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/chat"
)

type chatController struct {
	chatService service.ChatService
}

func NewChatController(chatService service.ChatService) *chatController {
	return &chatController{
		chatService: chatService,
	}
}

func (c *chatController) GetChatPage(e echo.Context) error {
	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return e.Redirect(http.StatusFound, "/login")
	}

	userId := e.Get("userId").(string)
	locationId := e.QueryParam("locationId")
	if locationId == "" {
		return e.HTML(http.StatusBadRequest, "Bad request")
	}

	messages, err := c.chatService.GetChat(userId, locationId)
	if err != nil {
		return e.HTML(http.StatusInternalServerError, "Internal server error")
	}

	return chat.Chat(locationId, messages).Render(e.Request().Context(), e.Response().Writer)
}

type SendChatForm struct {
	LocationId string `form:"locationId"`
	Content    string `form:"content"`
}

func (c *chatController) SendChat(e echo.Context) error {
	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return e.Redirect(http.StatusFound, "/login")
	}

	userId := e.Get("userId").(string)
	var form SendChatForm
	if err := e.Bind(&form); err != nil {
		return e.HTML(http.StatusBadRequest, "Bad request")
	}

	if form.LocationId == "" || form.Content == "" {
		return e.HTML(http.StatusBadRequest, "Bad request")
	}

	resp, err := c.chatService.SendChat(userId, form.LocationId, form.Content)
	if err != nil {
		return e.HTML(http.StatusInternalServerError, "Internal server error")
	}

	chatsHTML := ""
	for _, v := range resp {
		chatWriter := bytes.NewBufferString("")
		if v.Type == 0 {
			chat.TextMessage(v.Content, v.Sender).Render(e.Request().Context(), chatWriter)
			chatsHTML += chatWriter.String()
		} else if v.Type == 1 {
			chat.RecommendationMessage(form.LocationId, v.Content).Render(e.Request().Context(), chatWriter)
			chatsHTML += chatWriter.String()
		}
	}

	return e.HTML(http.StatusOK, chatsHTML)
}
