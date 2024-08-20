package cache

import (
	"context"
	"errors"
	"sync"

	"github.com/DimTur/chat-websocket-go/internal/domain"
	"github.com/google/uuid"
)

const chatDumpFileName = "chats.json"

type ChatsPool struct {
	sync.Mutex
	// key - ChatID
	pool map[domain.ID]*domain.Chat
}

func ChatCacheInit(ctx context.Context, wg *sync.WaitGroup) (*ChatsPool, error) {
	var chats = ChatsPool{
		pool: make(map[domain.ID]*domain.Chat),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(chatDumpFileName, chats.pool)
	}()

	if err := loadFromDump(chatDumpFileName, &chats.pool); err != nil {
		return nil, err
	}

	return &chats, nil
}

func (p *ChatsPool) AddChat(userIDs []domain.ID) domain.ID {
	chatID := domain.ID(uuid.New().String())
	nc := domain.Chat{
		UserIDs: userIDs,
		ChatID:  chatID,
	}

	p.Lock()
	p.pool[chatID] = &nc
	p.Unlock()

	return chatID
}

func (p *ChatsPool) AddMessage(chatID domain.ID, message domain.Message) error {
	p.Lock()
	defer p.Unlock()

	chat, ok := p.pool[chatID]
	if !ok {
		return errors.New("chat not found")
	}

	chat.Messages = append(chat.Messages, message)
	p.pool[chatID] = chat

	return nil
}

func (p *ChatsPool) DeleteMessage(chatID domain.ID, messageID domain.ID) error {
	p.Lock()
	defer p.Unlock()

	chat, ok := p.pool[chatID]
	if !ok {
		return errors.New("chat not found")
	}

	// TODO think
	for i, m := range chat.Messages {
		if m.MsgID == messageID {
			chat.Messages = append(chat.Messages[:i], chat.Messages[i+1:]...)
			break
		}
	}
	p.pool[chatID] = chat
	return nil
}

func (p *ChatsPool) UpdateMessage(chatID domain.ID, message domain.Message) error {
	p.Lock()
	defer p.Unlock()

	chat, ok := p.pool[chatID]
	if !ok {
		return errors.New("chat not found")
	}

	// TODO think
	for i, m := range chat.Messages {
		if m.MsgID == message.MsgID {
			// TODO надо проверить, тот же автор ли?
			chat.Messages[i] = message
			break
		}
	}
	p.pool[chatID] = chat
	return nil
}

func (p *ChatsPool) GetChatUsers(chatID domain.ID) ([]domain.ID, error) {
	p.Lock()
	defer p.Unlock()

	chat, ok := p.pool[chatID]
	if !ok {
		return nil, errors.New("chat not found")
	}
	return chat.UserIDs, nil
}
