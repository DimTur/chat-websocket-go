package domain

import (
	"encoding/json"
	"time"
)

type MsgType string
type ReqType string
type DeliveryType string
type ID string

const (
	MsgTypeAdd MsgType = "add"
	MsgTypeUpd MsgType = "upd"
	MsgTypeDel MsgType = "del"

	ReqTypeNewChat ReqType = "new_chat"
	ReqTypeNewMsg  ReqType = "new_msg"

	DeliveryTypeNewMsg  DeliveryType = "new_msg"
	DeliveryTypeNewChat DeliveryType = "new_chat"
	DeliveryTypeError   DeliveryType = "error"
)

// Chat and Message is storage data
type Chat struct {
	UserIDs  []ID      `json:"user_ids"`
	ChatID   ID        `json:"chat_id"`
	Messages []Message `json:"messages"`
}

type Message struct {
	MsgID  ID        `json:"msg_id"`
	Body   string    `json:"body"`
	TDate  time.Time `json:"t_date"`
	FromID ID        `json:"from_id"`
}

// Request
// Data for Type = ReqTypeNewChat is NewChatRequest
// Data for Type = ReqTypeNewMsg is MessageChatRequest
type Request struct {
	Type ReqType         `json:"type"`
	Data json.RawMessage `json:"data"`
}

type NewChatRequest struct {
	UserIDs []ID `json:"users"`
}

// MessageChatRequest can support messages types:
// MsgTypeAdd is new msg in chat
// MsgTypeUpd is changing existing msg in chat
// MsgTypeDel is deleting msg from chat
type MessageChatRequest struct {
	Msg    string  `json:"msg"`
	Type   MsgType `json:"type"`
	ChatID ID      `json:"chat_id"`
}

// Delivery is msg that go to a client
// Data for Type = DeliveryTypeNewChat is chatID string
// Data for Type = DeliveryTypeNewMsg is MessageChatDelivery
// Data for Type = DeliveryTypeError is ErrorResponse
type Delivery struct {
	Type DeliveryType `json:"type"`
	Data interface{}  `json:"data"`
}

// MessageChatDelivery can support messages types:
// MsgTypeAdd is new msg in chat
// MsgTypeUpd is changing existing msg in chat
// MsgTypeDel is deleting msg from chat
type MessageChatDelivery struct {
	Message
	Type   MsgType `json:"type"`
	ChatID ID      `json:"chat_id"`
}

// ErrorResponse can be supplemented with additional field "code"
// that client can easily handle
type ErrorResponse struct {
	Error string `json:"error"`
}
