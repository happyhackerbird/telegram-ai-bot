package controllers

import (
	"encoding/json"
	"example/plushie/plushie-bot/memory"
	"example/plushie/plushie-bot/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var api_key string
var AImodel = "pplx-70b-chat"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api_key = os.Getenv("AI_API_KEY")
}

func GetAIResponse(str string) string {
	response, err := queryAPI(str)
	if err != nil {
		return fmt.Sprintf("An error occured while querying the AI.")
	}
	var result model.Response
	err = json.Unmarshal(response, &result)
	if err != nil || len(result.Choices) == 0 {
		return fmt.Sprintf("An error occured while querying the AI.")
	}
	memory.AppendToHistory(result.Choices[0].Message.Content)
	return result.Choices[0].Message.Content
}

func queryAPI(str string) ([]byte, error) {
	messages, err := memory.CurrentMessageWithHistory(str)
	if err != nil {
		fmt.Printf("client: could not get history window: %s\n", err)
		return nil, err
	}
	jsonMsg, err := json.Marshal(messages)
	if err != nil {
		fmt.Printf("client: could not marshal json: %s\n", err)
		return nil, err
	}

	url := "https://api.perplexity.ai/chat/completions"
	payload := strings.NewReader("{\"model\":\"" + AImodel + "\",\"messages\":" + string(jsonMsg) + ",\"temperature\":1.1}")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: error reading response body: %s\n", err)
		return nil, err
	}

	return body, nil

}
