package model

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
)

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Message      AIMessage `json:"message"`
	FinishReason string    `json:"finish_reason"`
}

// Define the struct for completion details
type Completion struct {
	Choices []Choice `json:"choices"`
	Model   string   `json:"model"`
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Usage   Usage    `json:"usage"`
}

// Define the struct for usage details
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Define the struct for price details
type Price struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
	Total  float64 `json:"total"`
}

// Define the struct for word count
type Words struct {
	Input  int `json:"input"`
	Output int `json:"output"`
	Total  int `json:"total"`
}

// Define the main struct that includes all other structs
type Response struct {
	Data struct {
		Completion Completion `json:"completion"`
		Price      Price      `json:"price"`
		Words      Words      `json:"words"`
	} `json:"data"`
	Success bool `json:"success"`
}

type DBResponse struct {
	Code int `json:"code"`
	Data struct {
		InsertCount int      `json:"insertCount"`
		InsertId    []string `json:"insertIds"`
	}
}

type SemanticSearchResponse struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

type Error struct {
	Message string `json:"message"`
}

type Content struct {
	Distance float32 `json:"distance"`
	Text     string  `json:"message_text"`
}

func DecodeResponse(b []byte) (interface{}, error) {
	var r SemanticSearchResponse
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	switch r.Code {
	case 200:
		var c []Content
		err = json.Unmarshal(r.Data, &c)
		if err != nil {
			return nil, err
		}
		return c, nil

	default:
		var e Error
		err = json.Unmarshal(r.Data, &e)
		if err != nil {
			return nil, fmt.Errorf("cannot handle type: %s", r.Data)
		}
		return e, fmt.Errorf(e.Message)

	}
}

func GetMostRelevantMessages(c []Content, n int) []Content {
	slices.SortFunc(c, func(i Content, j Content) int {
		return -cmp.Compare(i.Distance, j.Distance)
	})
	return c[:n]
}

func GetResponseMessages(c []Content) []string {
	var messages []string
	for _, v := range c {
		fmt.Println(v.Distance)
		messages = append(messages, v.Text)
	}
	return messages
}
