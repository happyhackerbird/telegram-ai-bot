package model

import (
	"encoding/json"
	"fmt"
)

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Id      string `json:"id"`
	Model   string `json:"model"`
	Created int    `json:"created"`
	Usage   struct {
		Prompt_tokens     int `json:"prompt_tokens"`
		Completion_tokens int `json:"completion_tokens"`
		Total_tokens      int `json:"total_tokens"`
	}
	Object  string `json:"object"`
	Choices []struct {
		Index         int    `json:"index"`
		Finish_reason string `json:"finish_reason"`
		Message       AIMessage
		Delta         AIMessage
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

func DecodeMessage(b []byte) (interface{}, error) {
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

func GetResponseMessages(c []Content) []string {
	var messages []string
	for _, v := range c {
		messages = append(messages, v.Text)
	}
	return messages
}
