package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Connect() client.Client {
	ctx := context.Background()
	cl, err := client.NewClient(ctx, client.Config{
		Address: os.Getenv("DB_URL"),
		APIKey:  os.Getenv("DB_APITOKEN"),
	})
	if err != nil {
		log.Fatal("fail to connect to milvus", err.Error())
	}
	fmt.Println("Successfully connected to DB!")
	return cl
}
