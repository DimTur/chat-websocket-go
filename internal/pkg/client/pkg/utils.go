package utils

import (
	"bufio"
	"os"
	"strings"

	"github.com/DimTur/chat-websocket-go/internal/pkg/client/internal/models"
)

func GetUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func UInputUID() models.ID {
	reader := bufio.NewReader(os.Stdin)
	userID, _ := reader.ReadString('\n')
	userID = userID[:len(userID)-1]
	return models.ID(userID)
}

func UInputMsg(inputChan chan string) {
	reader := bufio.NewReader(os.Stdin)
	msg, _ := reader.ReadString('\n')
	inputChan <- msg
}
