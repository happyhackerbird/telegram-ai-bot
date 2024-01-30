package controllers

import (
	"encoding/json"
	memory "example/bot/telegram-ai-bot/database"
	"example/bot/telegram-ai-bot/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var api_key string
var AImodel = "pplx-70b-chat"
var temp = 1.1

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api_key = os.Getenv("AI_API_KEY")
}

func GetAIResponse(input string) string {
	// query the api
	response, err := queryAPI(input)
	if err != nil {
		return fmt.Sprintf("An error occured while querying the AI.")
	}
	// get json response
	var result model.Response
	err = json.Unmarshal(response, &result)
	if err != nil || len(result.Choices) == 0 {
		return fmt.Sprintf("An error occured while querying the AI.")
	}
	// append the response to the history if valid
	memory.AppendToHistory(result.Choices[0].Message.Content)
	return result.Choices[0].Message.Content
}

func queryAPI(input string) ([]byte, error) {
	// get a history window that fits the context length
	messages, err := memory.CurrentMessageWithHistory(input)
	if err != nil {
		fmt.Printf("client: could not get history window: %s\n", err)
		return nil, err
	}
	// fmt.Printf("messages: %v\n", messages)
	// format as json
	jsonMsg, err := json.Marshal(messages)
	if err != nil {
		fmt.Printf("client: could not marshal json: %s\n", err)
		return nil, err
	}

	// build the request and send it
	url := "https://api.perplexity.ai/chat/completions"
	s := strconv.FormatFloat(temp, 'f', -1, 64)
	payload := strings.NewReader("{\"model\":\"" + AImodel + "\",\"messages\":" + string(jsonMsg) + ",\"temperature\":\"" + s + "\"}")

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
