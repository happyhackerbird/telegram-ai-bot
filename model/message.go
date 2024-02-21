package model

type VectorizedMessage struct {
	// MessageID   int64  `json:"message_id" milvus:"message_id"`
	ChatID      int64  `json:"chat_id" milvus:"chat_id"`
	MessageText string `json:"message_text" milvus:"message_text"`
	// WordCount     int64     `json:"word_count" milvus:"word_count"`
	MessageVector []float32 `json:"message_vector" milvus:"message_vector"`
}