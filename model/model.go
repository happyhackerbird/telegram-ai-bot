package model

type Message struct {
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
		Message       Message
		Delta         Message
	}
}
