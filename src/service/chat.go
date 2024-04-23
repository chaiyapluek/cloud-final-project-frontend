package service

import (
	"encoding/base64"
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
		idx := strings.Index(pair, "=")
		key := pair[:idx]
		val := pair[idx+1:]
		if key != "preference" {
			continue
		}

		if val == "" {
			return true
		}
		m := map[string][]string{}
		b, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return true
		}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return true
		}
		if len(m) == 0 {
			return true
		}

		isEmpty := true
		for _, v := range m {
			if len(v) != 0 {
				isEmpty = false
				break
			}
		}
		return isEmpty
	}
	return true
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
