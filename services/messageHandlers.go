package services

import (
	"example/bot/telegram-ai-bot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StoreMessage(reply *tgbotapi.MessageConfig) {
	txt := reply.Text
	msg := &model.VectorizedMessage{
		MessageID:   b.GetCount(),
		ChatID:      reply.ChatID,
		MessageText: txt,
		// WordCount: getWordCount(txt),
		MessageVector: getVector(txt),
	}
	b.Store(msg)
}

func getVector(str string) []float32 {

}
