package service

import (
	"time"

	"github.com/DimTur/chat-websocket-go/internal/domain"
	"github.com/DimTur/chat-websocket-go/internal/repository/chatdb"
	"github.com/DimTur/chat-websocket-go/internal/service/pools"
	"github.com/google/uuid"
)

var chats chatdb.DB

func Init(chatdb chatdb.DB) {
	chats = chatdb
}

func NewMessage(msgReq domain.MessageChatRequest, fromID domain.ID) error {
	msg := domain.Message{
		MsgID:  domain.ID(uuid.New().String()),
		Body:   msgReq.Msg,
		TDate:  time.Now(),
		FromID: fromID,
	}

	if err := chats.AddMessage(msgReq.ChatID, msg); err != nil {
		return err
	}

	users, err := chats.GetChatUsers(msgReq.ChatID)
	if err != nil {
		return err
	}

	toDelivery := domain.MessageChatDelivery{
		Message: msg,
		Type:    msgReq.Type,
		ChatID:  msgReq.ChatID,
	}

	for _, userID := range users {
		if userID != fromID {
			pools.Users.Send(userID, toDelivery)
		}
	}

	return nil
}

func NewChat(userIDs []domain.ID) domain.ID {
	return chats.AddChat(userIDs)
}
