package database

import (
	"errors"
	"log"
	"strings"

	"example/bot/telegram-ai-bot/model"

	"github.com/pandodao/tokenizer-go"
)

var (
	systemText       string
	botInstruction   model.AIMessage
	systemTokenCount int
	globalHistory    []string
	answerLength     = 218  // number of tokens in medium length answer
	contextLength    = 4096 // LLM model context length
)

func SetInstruction(prompt string) {
	systemText = prompt
	botInstruction = model.AIMessage{
		Role:    "system",
		Content: systemText,
	}
	systemTokenCount = tokenizer.MustCalToken(systemText)
}

func CurrentMessageWithHistory(userMessage string) ([]model.AIMessage, error) {
	AppendToHistory(userMessage)
	err := getHistoryWindow()
	if err != nil {
		return nil, err
	}
	return getMsgObjects()
}

func getHistoryWindow() error {
	userMsgTokenCount := tokenizer.MustCalToken(globalHistory[len(globalHistory)-1])

	if userMsgTokenCount-answerLength > contextLength {
		return errors.New("user message too long")
	}
	tokenCount := tokenizer.MustCalToken(strings.Join(globalHistory, ".")) + systemTokenCount

	// use dequeue to remove the first element
	for tokenCount-answerLength > contextLength {
		log.Println("lenght of history:", len(globalHistory))
		tokenCount -= tokenizer.MustCalToken(globalHistory[0])
		globalHistory = globalHistory[1:]
	}
	log.Println("lenght of history after:", len(globalHistory))

	return nil
}

func getMsgObjects() ([]model.AIMessage, error) {
	// msgObjects := make([]model.AIMessage, len(globalHistory))
	msgObjects := []model.AIMessage{}
	msgObjects = append(msgObjects, botInstruction)

	user := true
	for _, msg := range globalHistory {
		msgObjects = append(msgObjects, stringToMessage(user, msg))
		user = !user
	}
	if user {
		return nil, errors.New("error trying to resize the message history")
	}
	return msgObjects, nil
}

func AppendToHistory(str string) {
	globalHistory = append(globalHistory, str)
	// log.Printf("new History: %v \n", globalHistory)
}

func stringToMessage(user bool, str string) model.AIMessage {
	if user {
		return model.AIMessage{
			Role:    "user",
			Content: str,
		}
	} else {
		return model.AIMessage{
			Role:    "assistant",
			Content: str,
		}
	}
}
