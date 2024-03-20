package model

type VectorizedMessage struct {
	ChatID      int64  `json:"chat_id" milvus:"chat_id"`
	MessageText string `json:"message_text" milvus:"message_text"`
	// WordCount     int64     `json:"word_count" milvus:"word_count"`
	MessageVector []float32 `json:"message_vector" milvus:"message_vector"`
}

type VectorizedProfile struct {
	ChatID   int64        `json:"chat_id" milvus:"name:chat_id"`
	Vector   []float32    `json:"vector" milvus:"name:vector"`
	Profiles JSONProfiles `json:"profiles" milvus:"name:profiles"`
}

type Profile struct {
	Name        string `json:"name" milvus:"name:name"`
	Instruction string `json:"instruction" milvus:"name:instruction"`
	AIModel     string `json:"ai_model" milvus:"name:ai_model"`
}

type JSONProfiles []Profile
