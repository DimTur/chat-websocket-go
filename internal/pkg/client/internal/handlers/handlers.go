package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DimTur/chat-websocket-go/internal/pkg/client/internal/models"
	utils "github.com/DimTur/chat-websocket-go/internal/pkg/client/pkg"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func ChatInputHandler(ctx context.Context, conn *websocket.Conn, chatID models.ID, userID models.ID) {
	var req models.Request
	req.Type = models.ReqTypeNewMsg

	var msgReq models.MessageChatRequest
	msgReq.ChatID = chatID
	msgReq.Type = models.MsgTypeAdd

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
			inputChan := make(chan string)
			go utils.UInputMsg(inputChan)

			input := <-inputChan
			if strings.TrimSpace(input) == "" {
				continue
			}

			msgReq.Msg = models.Message{
				MsgID:  models.ID(uuid.New().String()),
				Body:   input,
				TDate:  time.Now(),
				FromID: userID,
			}

			data, err := json.Marshal(msgReq)
			if err != nil {
				log.Println("Error serializing message:", err)
				return
			}
			req.Data = data

			if err := conn.WriteJSON(req); err != nil {
				log.Println("Error sending message:", err)
				return
			}
		}
	}
}

func CreateChat(conn *websocket.Conn, userID models.ID) (models.ID, error) {
	var req models.Request
	req.Type = models.ReqTypeNewChat

	var chatReq models.NewChatRequest
	chatReq.UserIDs = []models.ID{userID}
	data, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("error while serializing data: %w", err)
	}
	req.Data = data

	// Sending request to create a chat
	if err := conn.WriteJSON(req); err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	// Reading respinse from the server
	var resp models.Delivery
	if err := conn.ReadJSON(&resp); err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	chatID, ok := resp.Data.(string)
	if !ok {
		return "", fmt.Errorf("server error: not string")
	}

	return models.ID(chatID), nil
}

func HandleIncomingMsg(ctx context.Context, conn *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			// log.Printf("Message received: %s", message)
			var msg models.Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Error unmarshaling message:", err)
				continue
			}

			// Format msg and present it
			formattedMessage := formatMessage(msg)
			fmt.Println(formattedMessage)
		}
	}
}

func formatMessage(msg models.Message) string {
	formattedTime := msg.TDate.Format("02.01.2006 15:04:05")
	return fmt.Sprintf(">\n%s %s:\n%s>", msg.FromID, formattedTime, msg.Body)
}
