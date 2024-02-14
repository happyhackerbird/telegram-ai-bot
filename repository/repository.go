package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type Repository struct {
	Message Message
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Init() *Repository {
	// this goes to database package??
	ctx := context.Background()
	cl, err := client.NewClient(ctx, client.Config{
		Address: os.Getenv("DB_URL"),
		APIKey:  os.Getenv("DB_APITOKEN"),
	})
	if err != nil {
		log.Fatal("fail to connect to milvus", err.Error())
	}
	fmt.Println("Successfully connected to DB!")

	return &Repository{Message{cl}}

}
