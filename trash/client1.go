package main

// import (
// 	"bufio"
// 	"context"
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"strings"
// 	"syscall"
// 	"time"

// 	"github.com/DimTur/chat-websocket-go/internal/domain"
// 	"github.com/google/uuid"
// 	"github.com/gorilla/websocket"
// )

// func main() {
// 	userID := flag.String("userID", "", "User ID")
// 	flag.Parse()

// 	if *userID == "" {
// 		reader := bufio.NewReader(os.Stdin)
// 		fmt.Print("Введите ваш User ID: ")
// 		inputID, _ := reader.ReadString('\n')
// 		*userID = inputID[:len(inputID)-1] // удаляем символ новой строки
// 	}

// 	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

// 	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8001?userID="+*userID, nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer conn.Close()

// 	for {
// 		displayMenu()
// 		choice, err := getUserInput()
// 		if err != nil {
// 			log.Println("Error reading input:", err)
// 			continue
// 		}

// 		switch choice {
// 			case
// 		}
// 	}

// 	// Reading user ID from console
// 	reqUID := uInputUID()

// 	// Create a new chat
// 	chatID, err := createChat(conn, reqUID)
// 	if err != nil {
// 		log.Println("Error to create a chat", err)
// 		return
// 	}
// 	log.Println("Chat created with ID:", chatID)

// 	go func() {
// 		chatInputHandler(ctx, conn, chatID, reqUID)
// 	}()

// 	// Goroutine for reading messages from the server
// 	go handleIncomingMsg(ctx, conn)

// 	<-ctx.Done()
// }

// func chatInputHandler(ctx context.Context, conn *websocket.Conn, chatID domain.ID, userID domain.ID) {
// 	var req domain.Request
// 	req.Type = domain.ReqTypeNewMsg

// 	var msgReq domain.MessageChatRequest
// 	msgReq.ChatID = chatID
// 	msgReq.Type = domain.MsgTypeAdd

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log.Println("interrupt")

// 			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
// 			if err != nil {
// 				log.Println("write close:", err)
// 				return
// 			}
// 			return

// 		default:
// 			inputChan := make(chan string)
// 			go uInputMsg(inputChan)

// 			input := <-inputChan
// 			if strings.TrimSpace(input) == "" {
// 				continue
// 			}

// 			msgReq.Msg = domain.Message{
// 				MsgID:  domain.ID(uuid.New().String()),
// 				Body:   input,
// 				TDate:  time.Now(),
// 				FromID: userID,
// 			}

// 			data, err := json.Marshal(msgReq)
// 			if err != nil {
// 				log.Println("Error serializing message:", err)
// 				return
// 			}
// 			req.Data = data

// 			if err := conn.WriteJSON(req); err != nil {
// 				log.Println("Error sending message:", err)
// 				return
// 			}
// 		}
// 	}
// }

// func createChat(conn *websocket.Conn, userID domain.ID) (domain.ID, error) {
// 	var req domain.Request
// 	req.Type = domain.ReqTypeNewChat

// 	var chatReq domain.NewChatRequest
// 	chatReq.UserIDs = []domain.ID{userID}
// 	data, err := json.Marshal(chatReq)
// 	if err != nil {
// 		return "", fmt.Errorf("error while serializing data: %w", err)
// 	}
// 	req.Data = data

// 	// Sending request to create a chat
// 	if err := conn.WriteJSON(req); err != nil {
// 		return "", fmt.Errorf("error sending request: %w", err)
// 	}

// 	// Reading respinse from the server
// 	var resp domain.Delivery
// 	if err := conn.ReadJSON(&resp); err != nil {
// 		return "", fmt.Errorf("error reading response: %w", err)
// 	}

// 	chatID, ok := resp.Data.(string)
// 	if !ok {
// 		return "", fmt.Errorf("server error: not string")
// 	}

// 	return domain.ID(chatID), nil
// }

// func handleIncomingMsg(ctx context.Context, conn *websocket.Conn) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			_, message, err := conn.ReadMessage()
// 			if err != nil {
// 				log.Println("Error reading message:", err)
// 				return
// 			}
// 			// log.Printf("Message received: %s", message)
// 			var msg domain.Message
// 			if err := json.Unmarshal(message, &msg); err != nil {
// 				log.Println("Error unmarshaling message:", err)
// 				continue
// 			}

// 			// Format msg and present it
// 			formattedMessage := formatMessage(msg)
// 			fmt.Println(formattedMessage)
// 		}
// 	}
// }

// func uInputUID() domain.ID {
// 	reader := bufio.NewReader(os.Stdin)
// 	userID, _ := reader.ReadString('\n')
// 	userID = userID[:len(userID)-1]
// 	return domain.ID(userID)
// }

// func uInputMsg(inputChan chan string) {
// 	reader := bufio.NewReader(os.Stdin)
// 	msg, _ := reader.ReadString('\n')
// 	inputChan <- msg
// }

// func formatMessage(msg domain.Message) string {
// 	formattedTime := msg.TDate.Format("02.01.2006 15:04:05")
// 	return fmt.Sprintf(">\n%s %s:\n%s>", msg.FromID, formattedTime, msg.Body)
// }

// func displayMenu() {
// 	fmt.Println("1. Создать новый чат с другим пользователем")
// 	fmt.Println("2. Войти в чат с пользователем")
// 	fmt.Println("Введите ваш выбор (для выхода введите exit или нажмите ctrl+C):")
// }

// func getUserInput() (string, error) {
// 	reader := bufio.NewReader(os.Stdin)
// 	input, err := reader.ReadString('\n')
// 	if err != nil {
// 		return "", err
// 	}
// 	return strings.TrimSpace(input), nil
// }
