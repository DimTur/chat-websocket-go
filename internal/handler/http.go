package handler

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DimTur/chat-websocket-go/internal/domain"
	"github.com/DimTur/chat-websocket-go/internal/pkg/authclient"
	"github.com/gorilla/websocket"
)

const (
	HeaderAuthorization = "Authorization"
	HeaderUserID        = "User-ID"
	HeaderUserRole      = "User-Role"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Minute,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	WriteBufferPool:  &sync.Pool{},
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func HandleHTTPReq(resp http.ResponseWriter, req *http.Request) {
	defer func() {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Add("Access-Control-Allow-Methods", "GET")
		resp.Header().Add("Access-Control-Allow-Methods", "OPTIONS")
	}()

	token := req.Header.Get(HeaderAuthorization)

	if token == "" {
		resp.WriteHeader(http.StatusUnauthorized)
		log.Println("Get request", req.Method, token, "error", http.StatusUnauthorized)
		return
	}

	userID, valid := authclient.ValidateToken(token)
	if !valid {
		resp.WriteHeader(http.StatusUnauthorized)
		log.Println("Get request", req.Method, token, "error", http.StatusUnauthorized)
		return
	}

	log.Println("userID", userID)

	// req.Header.Set(HeaderUserID, userID)

	// обновление соединения до WebSocket
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		// Upgrade сам вставляет статус код в респонс
		log.Println(err)
		return
	}
	HandleWsConn(conn, domain.ID(userID))
}

// func HandleHTTPReq(resp http.ResponseWriter, req *http.Request) {
// 	defer func() {
// 		resp.Header().Set("Access-Control-Allow-Origin", "*")
// 		resp.Header().Add("Access-Control-Allow-Methods", "GET")
// 		resp.Header().Add("Access-Control-Allow-Methods", "OPTIONS")
// 	}()

// 	userID := req.URL.Query().Get("userID")

// 	if userID == "" {
// 		resp.WriteHeader(http.StatusBadRequest)
// 		log.Println("Get request", req.Method, "missing userID", "error", http.StatusBadRequest)
// 		return
// 	}

// 	log.Println("userID", userID)

// 	// req.Header.Set(HeaderUserID, userID.Hex())
// 	// req.Header.Set(HeaderUserID, userID)

// 	// обновление соединения до WebSocket
// 	conn, err := upgrader.Upgrade(resp, req, nil)
// 	if err != nil {
// 		// Upgrade сам вставляет статус код в респонс
// 		log.Println(err)
// 		return
// 	}
// 	HandleWsConn(conn, domain.ID(userID))
// }
