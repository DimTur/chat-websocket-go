package chatdb

import "github.com/DimTur/chat-websocket-go/internal/domain"

type DB interface {
	AddMessage(chatID domain.ID, message domain.Message) error
	DeleteMessage(chatID domain.ID, messageID domain.ID) error
	UpdateMessage(chatID domain.ID, message domain.Message) error
	GetChatUsers(chatID domain.ID) ([]domain.ID, error)
	AddChat(userIDs []domain.ID) domain.ID
}
