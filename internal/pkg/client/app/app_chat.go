package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/DimTur/chat-websocket-go/internal/pkg/client/internal/config"
	"github.com/DimTur/chat-websocket-go/internal/pkg/client/internal/menu"
	"github.com/gorilla/websocket"
)

func Run() {
	cfg := config.GetConfig()

	headers := http.Header{}
	headers.Set("Authorization", cfg.AuthorizationToken)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8001", headers)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	var userData map[string]string
	if err := conn.ReadJSON(&userData); err != nil {
		log.Fatal("Error reading initial userID:", err)
	}

	userID := userData["userID"]
	if userID == "" {
		log.Fatal("Received empty userID from server")
	}

	fmt.Printf("Received userID: %s\n", userID)

	go menu.ChatMenu(ctx, conn, userID)

	<-ctx.Done()
	fmt.Println("Exiting...")
}
