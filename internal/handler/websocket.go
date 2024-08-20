package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DimTur/chat-websocket-go/internal/domain"
	"github.com/DimTur/chat-websocket-go/internal/service"
	"github.com/DimTur/chat-websocket-go/internal/service/pools"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 1 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func HandleWsConn(conn *websocket.Conn, userID domain.ID) {
	defer func() {
		// closing the user channel and ending write goroutine
		if pools.Users.Delete(userID) {
			conn.Close()
		}
	}()

	ch := pools.Users.New(userID)

	// write to conn from channel
	go func() {
		tiker := time.NewTicker(pingPeriod)
		defer func() {
			tiker.Stop()
			if pools.Users.Delete(userID) {
				conn.Close()
			}
		}()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Println("channel for " + userID + " closed")
				}
				log.Println("SEND to ", userID, msg)

				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteJSON(msg); err != nil {
					handleWsError(err, userID)
					return
				}
			case <-tiker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					handleWsError(err, userID)
					return
				}
			}
		}
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	// read from conn
	for {
		typ, message, err := conn.ReadMessage()
		if err != nil {
			handleWsError(err, userID)
			return
		}
		log.Println("GOT from", userID, typ, string(message))

		switch typ {
		case websocket.TextMessage, websocket.BinaryMessage:
			var req domain.Request
			if err = json.Unmarshal(message, &req); err != nil {
				sendErrorResp(userID, err)
				continue
			}

			switch req.Type {
			case domain.ReqTypeNewChat:
				var newChatReq domain.NewChatRequest
				if err = json.Unmarshal(req.Data, &newChatReq); err != nil {
					sendErrorResp(userID, err)
					continue
				}
				chatID := service.NewChat(append(newChatReq.UserIDs, userID))
				sendResp(userID, domain.DeliveryTypeNewChat, chatID)

			case domain.ReqTypeNewMsg:
				var msg domain.MessageChatRequest
				if err = json.Unmarshal(req.Data, &msg); err != nil {
					sendErrorResp(userID, err)
					continue
				}

				switch msg.Type {
				case domain.MsgTypeAdd:
					if err := service.NewMessage(msg, userID); err != nil {
						sendErrorResp(userID, err)
						continue
					}
				}
			}

		case websocket.CloseMessage:
			return
		}
	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}

func handleWsError(err error, userID domain.ID) {
	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
		log.Println("websocket session closed by client", userID)
	default:
		log.Println("error websocket message", err.Error(), "for", userID)
	}
}

func sendErrorResp(userID domain.ID, err error) {
	sendResp(userID, domain.DeliveryTypeError, domain.ErrorResponse{Error: err.Error()})
}

func sendResp(userID domain.ID, typ domain.DeliveryType, data interface{}) {
	var resp domain.Delivery
	resp.Type = typ
	resp.Data = data
	pools.Users.Send(userID, resp)
}
