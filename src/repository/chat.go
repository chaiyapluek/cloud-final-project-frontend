package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"dev.chaiyapluek.cloud.final.frontend/src/entity"
)

type ChatRepository interface {
	GetChat(userId, locationId string) ([]*entity.Chat, error)
	SendChat(userId, locationId, content string) ([]*entity.Chat, error)
}

type chatRepository struct {
	backendURL string
	apiKey     string
}

func NewChatRepository(backendURL string, apiKey string) ChatRepository {
	return &chatRepository{
		backendURL: backendURL,
		apiKey:     apiKey,
	}
}

func (r *chatRepository) GetChat(userId, locationId string) ([]*entity.Chat, error) {
	endpoint := "/users/" + userId + "/chats?locationId=" + locationId

	httpReq, err := http.NewRequest("GET", r.backendURL+endpoint, nil)
	if err != nil {
		log.Println("chat repository get "+endpoint+" error", err)
		return nil, err
	}
	httpReq.Header.Set("X-API-KEY", r.apiKey)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("chat repository get "+endpoint+" error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("chat repository get "+endpoint+" status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("chat repository get "+endpoint+" decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var chatResp entity.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		log.Println("chat repository get "+endpoint+" decode ChatResponse error", err)
		return nil, err
	}

	return chatResp.Data, nil
}

func (r *chatRepository) SendChat(userId, locationId, content string) ([]*entity.Chat, error) {
	endpoint := "/chats"

	req := entity.SendChatRequest{
		UserId:     userId,
		LocationId: locationId,
		Content:    content,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Println("chat repository send marshal error", err)
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", r.backendURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("chat repository send "+endpoint+" error", err)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", r.apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("chat repository send "+endpoint+" error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("chat repository send "+endpoint+" status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("chat repository send "+endpoint+" decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var chatResp entity.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		log.Println("chat repository send "+endpoint+" decode ChatResponse error", err)
		return nil, err
	}

	return chatResp.Data, nil
}
