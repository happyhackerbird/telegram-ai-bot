package controllers

import (
	"encoding/json"
	"example/bot/telegram-ai-bot/database"
	"example/bot/telegram-ai-bot/model"
	"example/bot/telegram-ai-bot/services"
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
var AIModel string
var systemText string
var temp = 1.1
var historyPrompt = "Below are some of the most relevant interactions from this chat. \n %v \n Use the information to respond to this message: %v"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api_key = os.Getenv("AI_API_KEY")
}

func SetModel(m string) {
	AIModel = m
}

func SetInstruction(prompt string) {
	systemText = prompt
}

func GetAIResponse(chatID int64, input string) string {
	// query the api
	response, err := queryAPI(chatID, input)
	if err != nil {
		return ("An error occured while querying the AI.")
	}
	// get json response
	var result model.Response
	err = json.Unmarshal(response, &result)
	if err != nil {
		log.Println("client: error querying the AI: ", err)
		return ("An error occured while querying the AI.")
	} else if len(result.Choices) == 0 {
		log.Println("client: error querying the AI: received empty response")
		log.Println(result)

		return ("An error occured while querying the AI.")

	}
	// append the response to the history if valid
	database.AppendToHistory(result.Choices[0].Message.Content)
	return result.Choices[0].Message.Content
}

func queryAPI(chatID int64, input string) ([]byte, error) {
	// get context
	context, err := services.GetContext(chatID, input)
	if err != nil {
		log.Printf("client: error getting context: %s\n", err)
	}
	fmt.Println("\nContext:\n", context)

	var messages []model.AIMessage
	botInstruction := model.AIMessage{
		Role:    "system",
		Content: systemText,
	}
	questionWithContext := model.AIMessage{
		Role:    "user",
		Content: fmt.Sprintf(historyPrompt, context, input),
	}
	messages = append(messages, botInstruction)
	messages = append(messages, questionWithContext)

	// format as json
	jsonMsg, err := json.Marshal(messages)
	if err != nil {
		log.Printf("client: could not marshal json: %s\n", err)
		return nil, err
	}

	s := string('"') + strconv.FormatFloat(temp, 'f', -1, 64) + string('"')
	// build the request and send it
	var url string
	if AIModel == "gpt-4-turbo-preview" {
		url = "https://api.openai.com/v1/chat/completions"
		api_key = "Bearer " + os.Getenv("OPENAI_API_KEY")
		s = strconv.FormatFloat(0.6, 'f', -1, 64) // openai takes between 0 and 1 as temp

	} else {
		url = "https://api.perplexity.ai/chat/completions"
	}
	payload := strings.NewReader("{\"model\":\"" + AIModel + "\",\"messages\":" + string(jsonMsg) + ",\"temperature\": " + s + "}")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Printf("client: could not create request: %s\n", err)
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("client: error making http request: %s\n", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("client: error reading response body: %s\n", err)
		return nil, err
	}
	return body, nil

}
