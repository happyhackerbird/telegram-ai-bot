package services

import (
	"context"
	"example/bot/telegram-ai-bot/model"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
)

func StoreMessage(reply *tgbotapi.MessageConfig) {
	txt := reply.Text
	if v, err := getVector(txt); err == nil {
		msg := &model.VectorizedMessage{
			MessageID:   b.GetCount(),
			ChatID:      reply.ChatID,
			MessageText: txt,
			// WordCount: getWordCount(txt),
			MessageVector: v,
		}
		b.Store(msg) // error at this level
	}
}

func getVector(str string) ([]float32, error) {
	fmt.Println("Getting vector embedding for message ... ")
	client := openai.NewClient(os.Getenv("OPENAI_API"))

	// Create an EmbeddingRequest for the user query
	queryReq := openai.EmbeddingRequest{
		Input: []string{str},
		Model: openai.SmallEmbedding3,
	}

	// Create an embedding for the user query
	queryResponse, err := client.CreateEmbeddings(context.Background(), queryReq)
	if err != nil {
		log.Println("Error creating query embedding:", err)
		return nil, err
	}
	fmt.Println("Dimension of embeddings vector:", len(queryResponse.Data[0].Embedding))
	return queryResponse.Data[0].Embedding, nil

}
