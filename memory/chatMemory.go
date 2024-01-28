package memory

import (
	"errors"
	"example/plushie/plushie-bot/model"
	"fmt"
	"strings"

	"github.com/pandodao/tokenizer-go"
)

var systemText = "have a conversation with me about plushies and their lives & characters. be imaginative and creative. do not give precise or factual answers. do not invent new characters."
var botInstruction = model.Message{
	Role:    "system",
	Content: systemText,
}
var globalHistory = []string{systemText}
var answerLength = 218
var contextLength = 4096

// []model.Message{botInstruction}

func CurrentMessageWithHistory(userMessage string) ([]model.Message, error) {
	AppendToHistory(userMessage)
	err := getHistoryWindow()
	if err != nil {
		return nil, err
	}
	return getMsgObjects()
	// appendToHistory(str)
	// return GetAIResponse(globalHistory)
}

func getHistoryWindow() error {
	userMsgTokenCount := tokenizer.MustCalToken(globalHistory[len(globalHistory)-1])

	if userMsgTokenCount-answerLength > contextLength {
		return errors.New("User message too long")
	}
	tokenCount := tokenizer.MustCalToken(strings.Join(globalHistory, "."))
	// messages := []model.Message{
	// 	botInstruction, userMessage,
	// }

	// use queue to remove the first element
	for tokenCount-answerLength > contextLength {
		tokenCount -= tokenizer.MustCalToken(globalHistory[0])
		globalHistory = globalHistory[1:]
	}
	return nil
}

func getMsgObjects() ([]model.Message, error) {
	// msgObjects := make([]model.Message, len(globalHistory))
	msgObjects := []model.Message{}
	msgObjects = append(msgObjects, botInstruction)

	user := true
	for _, msg := range globalHistory[1:] {
		msgObjects = append(msgObjects, stringToMessage(user, msg))
		user = !user
	}
	if user {
		return nil, errors.New("Error trying to resize the message history")
	}
	return msgObjects, nil
}

func AppendToHistory(str string) {
	globalHistory = append(globalHistory, str)
	fmt.Printf("new History: %v \n", globalHistory)
}

func stringToMessage(user bool, str string) model.Message {
	if user {
		return model.Message{
			Role:    "user",
			Content: str,
		}
	} else {
		return model.Message{
			Role:    "assistant",
			Content: str,
		}
	}
}
