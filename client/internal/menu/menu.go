package menu

import (
	"context"
	"fmt"
	"log"

	"client/internal/handlers"
	"client/internal/models"
	utils "client/pkg"

	"github.com/gorilla/websocket"
)

func ChatMenu(ctx context.Context, conn *websocket.Conn, userID string) {
	for {
		displayMenu()
		choice, err := utils.GetUserInput()
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		switch choice {
		case "1":
			// Create a new chat
			fmt.Print("Введите ID пользователя, с которым хотите создать чат\nили введите return для выхода в предыдущее меню.\n")
			reqUID := utils.UInputUID()

			if reqUID == "return" {
				fmt.Println("Возврат в главное меню...")
				continue
			}

			chatID, err := handlers.CreateChat(conn, reqUID)
			if err != nil {
				log.Println("Error creating chat:", err)
				continue
			}
			log.Println("Chat created with ID:", chatID)
			continue

		case "2":
			// Entering to chat
			fmt.Print("Введите ID чата, в котором хотите продолжить общение\n или введите return для выхода в предыдущее меню.\n")
			chatID := utils.UInputUID()

			if chatID == "return" {
				fmt.Println("Возврат в главное меню...")
				continue
			}

			log.Println("Entering chat with ID:", chatID)

			// Goroutine for reading messages from the server
			go handlers.HandleIncomingMsg(ctx, conn)

			go handlers.ChatInputHandler(ctx, conn, chatID, models.ID(userID))

			<-ctx.Done()
			fmt.Println("Exiting...")
			return

		case "exit":
			fmt.Println("Exiting...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error closing connection:", err)
			}
			return

		default:
			<-ctx.Done()
			fmt.Println("Некорректный выбор, попробуйте снова.")
		}
	}
}

func displayMenu() {
	fmt.Println("")
	fmt.Println("1. Создать новый чат с другим пользователем")
	fmt.Println("2. Войти в чат с пользователем")
	fmt.Println("Введите ваш выбор (для выхода введите exit или нажмите ctrl+C):")
	fmt.Println("")
}
