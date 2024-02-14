package repository

import (
	"encoding/json"
	"example/bot/telegram-ai-bot/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

var COLLECTIONNAME = "Messages"

type Message struct {
	db client.Client
	// text string
}

func (m *Message) Store(row *model.VectorizedMessage) error {
	jsonBytes, err := json.Marshal(row)
	if err != nil {
		log.Printf("Error marshaling struct: %v", err)
	}
	payload := strings.NewReader(`{"collectionName":"` + COLLECTIONNAME + `", "data":` + string(jsonBytes) + "}")
	url := fmt.Sprintf("%s/v1/vector/insert", os.Getenv("DB_URL"))
	api_key := "Bearer " + os.Getenv("DB_APITOKEN")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Printf("client: could not create request: %s\n", err)
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("client: error making http request: %s\n", err)
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("client: error reading response body: %s\n", err)
		return err
	}
	fmt.Println(string(body))

	// log.Printf("Inserted %v\n", result)
	return nil
}
