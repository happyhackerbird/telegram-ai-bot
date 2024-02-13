package repository

import (
	"context"
	"example/bot/telegram-ai-bot/model"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

var COLLECTIONNAME = "Messages"

type Message struct {
	db client.Client
	// text string
}

func (m *Message) Store(row *model.VectorizedMessage) error {
	result, err := m.db.InsertRows(
		context.Background(),
		COLLECTIONNAME,
		"",
		[]interface{}{row},
	)
	if err != nil {
		return err

	}
	log.Printf("Inserted %v\n", result)
	return nil
}
