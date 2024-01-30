package memory

import (
	"errors"
	"example/bot/telegram-ai-bot/model"
	"log"
	"strings"

	"github.com/pandodao/tokenizer-go"
)

var (
	systemText     = "have a conversation with me about plushies and their lives & characters. be imaginative and creative. do not give precise or factual answers. do not invent new characters."
	botInstruction = model.Message{
		Role:    "system",
		Content: systemText,
	}
	globalHistory = []string{systemText}
	answerLength  = 218  // number of tokens in medium length answer
	contextLength = 4096 // LLM model context length
)

func CurrentMessageWithHistory(userMessage string) ([]model.Message, error) {
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
	tokenCount := tokenizer.MustCalToken(strings.Join(globalHistory, "."))

	// use dequeue to remove the first element
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
		return nil, errors.New("error trying to resize the message history")
	}
	return msgObjects, nil
}

func AppendToHistory(str string) {
	globalHistory = append(globalHistory, str)
	log.Printf("new History: %v \n", globalHistory)
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
