package main

import (
	"context"
	"encoding/json"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimTur/chat-websocket-go/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	userID := uuid.New().String() // TODO my UID

	var req domain.Request
	req.Type = domain.ReqTypeNewChat

	var chatReq domain.NewChatRequest
	chatReq.UserIDs = []string{uuid.New().String()} // TODO send friend id
	data, err := json.Marshal(chatReq)
	if err != nil {
		return
	}
	req.Data = data

	if err := conn.WriteJSON(req); err != nil {
		return
	}

	var resp domain.Delivery
	if err := conn.ReadJSON(&resp); err != nil {
		log.Println("read:", err)
		return
	}
	chatID, ok := resp.Data.(string)
	if !ok {
		log.Println("server error: not string")
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: %s", message)
			}
		}
	}()

	req.Type = domain.ReqTypeNewChat

	var chatReq2 domain.MessageChatRequest
	chatReq2.ChatID = chatID
	chatReq2.Type = domain.ReqTypeNewMsg
	for {
		select {
		case <-ctx.Done():
			log.Println("interrupt")

			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		default:
			chatReq2.Msg = domain.Message{
				MsgID:  uuid.New().String(),
				Body:   "Hi",
				TDate:  time.Now(),
				FromID: userID,
			}
			data, err := json.Marshal(chatReq2)
			if err != nil {
				return
			}
			req.Data = data

			if err := conn.WriteJSON(req); err != nil {
				return
			}
			time.Sleep(time.Second)
		}
	}
}
