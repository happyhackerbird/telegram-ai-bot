package database

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var cl client.Client
var COLLECTION_NAME = "Messages"

func Migrate() {
	cl = Connect()

	// delete collection if exists
	ctx := context.Background()
	has, err := cl.HasCollection(ctx, COLLECTION_NAME)
	if err != nil {
		log.Fatal("fail to check whether collection exists", err.Error())
	}
	if has {
		cl.DropCollection(ctx, COLLECTION_NAME)
		fmt.Println("Collection dropped")
	}

	// create a collection
	fmt.Println("Creating collection")
	schema := &entity.Schema{
		CollectionName: COLLECTION_NAME,
		Description:    "AI Bot Chat History",
		Fields: []*entity.Field{
			{
				Name:       "message_id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:       "chat_id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: false,
				AutoID:     false,
			},
			{
				Name:     "message_text",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "12560",
				},
			},

			// {
			// 	Name:       "word_count",
			// 	DataType:   entity.FieldTypeInt64,
			// 	PrimaryKey: false,
			// 	AutoID:     false,
			// },
			{
				Name:     "message_vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "1536",
				},
			},
		},
	}
	// entity.NewSchema().WithName(collectionName).WithDescription("AI Bot Chat History").
	// 	WithField(entity.NewField().WithName("message_id").WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithDescription("primary id")).
	// 	WithField(entity.NewField().WithName("chat_id").WithDataType(entity.FieldTypeInt64).WithDescription("chat identifier")).
	// 	WithField(entity.NewField().WithName("message_text").WithDataType(entity.FieldTypeVarChar).WithDescription("message text").WithType).
	// 	WithField(entity.NewField().WithName("message_vector").WithDataType(entity.FieldTypeFloatVector).WithDim(128))

	err = cl.CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		log.Fatal("failed to create collection", err.Error())
	}
	fmt.Println("Successfully created collection!")

}
