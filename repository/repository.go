package repository

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type Repository struct {
	Message Message
	Profile Profile
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Init(cl client.Client) *Repository {
	return &Repository{Message{cl},
		Profile{cl}}

}
