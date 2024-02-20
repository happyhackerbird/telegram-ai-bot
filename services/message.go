package services

import (
	"context"
	"example/bot/telegram-ai-bot/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
)

var n = 10
var COLLECTION_NAME = "Messages"

func StoreMessage(reply *tgbotapi.MessageConfig) {
	txt := reply.Text
	if v, err := getVector(txt); err == nil {
		msg := &model.VectorizedMessage{
			ChatID:      reply.ChatID,
			MessageText: txt,
			// WordCount: getWordCount(txt),
			MessageVector: v,
		}
		b.Store(msg)
	}
}

func getVector(str string) ([]float32, error) {
	// fmt.Println("Getting vector embedding for message ... ")
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
	// fmt.Println("Dimension of embeddings vector:", len(queryResponse.Data[0].Embedding))
	return queryResponse.Data[0].Embedding, nil

}

// so should this be here, or split up with repository or database ?? (repository -> through b.repository )
func GetContext(chatID int64, input string) (string, error) {
	fmt.Println("Getting context for message ... ")
	// get vector embedding for user input
	inputVector, err := getVector(input)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/v1/vector/search", os.Getenv("DB_URL"))
	api_key := "Bearer " + os.Getenv("DB_APITOKEN")
	payload := strings.NewReader(`{"collectionName": "` + COLLECTION_NAME + `", "outputFields": ["message_text"], "vector":` + toString(inputVector) + fmt.Sprintf(`, "filter": "chat_id in [%v]", "limit": %v}`, chatID, n))
	// fmt.Println("Payload: ", payload)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	response, err := model.DecodeMessage(body)
	if err != nil {
		return "", err
	}

	return strings.Join(model.GetResponseMessages(response.([]model.Content)), "\n"), nil
}

func toString(v []float32) string {
	return strings.Join(strings.Split(fmt.Sprint(v), " "), ", ")
}
