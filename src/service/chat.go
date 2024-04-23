package service

import (
	"encoding/json"
	"strings"

	"dev.chaiyapluek.cloud.final.frontend/src/repository"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/chat"
)

type ChatService interface {
	GetChat(userId, locationId string) ([]*chat.Message, error)
	SendChat(userId, locationId, content string) ([]*chat.Message, error)
}

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func isEmptyTypeOneMessage(content string) bool {
	strs := strings.Split(content, "?")
	if len(strs) != 2 {
		return false
	}

	query := strs[1]
	if query == "" {
		return true
	}

	q := strings.Split(query, "&")
	for _, pair := range q {
		p := strings.Split(pair, "=")
		if p[0] != "preference=" {
			continue
		}

		if p[1] == "" {
			return true
		}

		m := map[string][]string{}
		json.Unmarshal([]byte(p[1]), &m)
		if len(m) == 0 {
			return true
		}

		for _, v := range m {
			if len(v) == 0 {
				break
			}
		}

		return true
	}

	return false
}

func (s *chatService) GetChat(userId, locationId string) ([]*chat.Message, error) {
	resp, err := s.chatRepo.GetChat(userId, locationId)
	if err != nil {
		return nil, err
	}

	messages := []*chat.Message{}
	for _, v := range resp {
		if v.Type == 1 && isEmptyTypeOneMessage(v.Content) {
			continue
		}
		messages = append(messages, &chat.Message{
			Type:    v.Type,
			Sender:  v.Sender,
			Content: v.Content,
		})
	}

	return messages, nil
}

func (s *chatService) SendChat(userId, locationId, content string) ([]*chat.Message, error) {
	resp, err := s.chatRepo.SendChat(userId, locationId, content)
	if err != nil {
		return nil, err
	}

	messages := []*chat.Message{}
	for _, v := range resp {
		if v.Type == 1 && isEmptyTypeOneMessage(v.Content) {
			continue
		}
		messages = append(messages, &chat.Message{
			Type:    v.Type,
			Sender:  v.Sender,
			Content: v.Content,
		})
	}

	return messages, nil
}
