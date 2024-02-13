package repository

import (
	"example/bot/telegram-ai-bot/model"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type Message struct {
	db client.Client
	// text string
}

func (m *Message) Store(m *model.Message) error {
	_, err := m.
	if err != nil {
		return err
	}
	return nil
}
