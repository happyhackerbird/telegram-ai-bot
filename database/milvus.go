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

var cl client.Client

func Connect() client.Client {
	ctx := context.Background()
	c, err := client.NewClient(ctx, client.Config{
		Address: os.Getenv("DB_URL"),
		APIKey:  os.Getenv("DB_APITOKEN"),
	})
	if err != nil {
		log.Fatal("fail to connect to milvus", err.Error())
	}
	fmt.Println("Successfully connected to DB!")
	cl = c
	return c
}

func LoadIndex() {
	// 4. Load collection
	loadCollErr := cl.LoadCollection(context.Background(), COLLECTION_NAME, false)

	if loadCollErr != nil {
		log.Fatal("Failed to load collection:", loadCollErr.Error())
	}

	// 5. Get load progress
	_, err := cl.GetLoadingProgress(context.Background(), COLLECTION_NAME, nil)

	if err != nil {
		log.Fatal("Failed to get loading progress:", err.Error())
	}

	fmt.Println("Index loaded successfully")
}
