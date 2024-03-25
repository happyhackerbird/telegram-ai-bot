package repository

import (
	"context"
	"encoding/json"
	"example/bot/telegram-ai-bot/model"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var COLLECTION_NAME_PROFILE = "Profiles"

type Profile struct {
	db client.Client
}

func (p *Profile) Store(profiles *model.VectorizedProfile) error {
	fmt.Println("Storing profile in database ... ")

	conn := p.db

	// schema := getProfileSchema()

	// var profilesJSON []byte
	// var err error
	// if profiles != nil {
	// 	profilesJSON, err = json.Marshal(profiles.Profiles)
	// 	if err != nil {
	// 		log.Printf("Error marshaling struct: %v", err)
	// 	}
	// } else {
	// 	log.Println("Empty profiles") // shouldnt happen
	// 	return errors.New("Empty profiles")
	// }

	// // Prepare the data for upsert
	// row := struct {
	// 	ChatID   int64     `json:"name" milvus:"name:name"`
	// 	Vector   []float32 `json:"vector" milvus:"name:vector"`
	// 	Profiles string    `json:"profiles" milvus:"name:profiles"`
	// }{
	// 	ChatID:   profiles.ChatID,
	// 	Vector:   profiles.Vector,
	// 	Profiles: string(profilesJSON),
	// }
	// columns, _ := entity.AnyToColumns([]interface{}{row}, schema)

	// res, err := conn.Upsert(context.Background(), COLLECTION_NAME_PROFILE, "", columns...)
	// if err != nil {
	// 	log.Printf("Failed to upsert rows: ", err.Error())
	// }
	// fmt.Println("Upserted ", res, " rows")

	jsonBytes, err := json.Marshal(profiles.Profiles)
	if err != nil {
		log.Printf("Error marshaling struct: %v", err)
		return err
	}
	profilesJSONObject := map[string]interface{}{
		"profiles": json.RawMessage(jsonBytes),
	}

	wrappedProfilesJSON, err := json.Marshal(profilesJSONObject)
	if err != nil {
		log.Printf("Failed to marshal wrapped profiles:", err.Error())
		return err
	}

	// Prepare the row for insertion
	row := map[string]interface{}{
		"chat_id":  profiles.ChatID,
		"vector":   profiles.Vector,
		"profiles": wrappedProfilesJSON,
	}

	// Insert the data
	_, err = conn.InsertRows(context.Background(), COLLECTION_NAME_PROFILE, "", []interface{}{row})
	if err != nil {
		log.Printf("Failed to insert rows:", err.Error())
		return err
	}

	// payload := strings.NewReader(`{"collectionName": "` + COLLECTION_NAME_PROFILE + `", "data":` + string(jsonBytes) + "}")
	// // fmt.Println(payload)
	// url := fmt.Sprintf("%s/v1/vector/insert", os.Getenv("DB_URL"))
	// api_key := "Bearer " + os.Getenv("DB_APITOKEN")

	// req, err := http.NewRequest("POST", url, payload)
	// if err != nil {
	// 	log.Printf("client: could not create request: %s\n", err)
	// 	return err
	// }

	// req.Header.Add("accept", "application/json")
	// req.Header.Add("content-type", "application/json")
	// req.Header.Add("authorization", api_key)

	// res, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Printf("client: error making http request: %s\n", err)
	// 	return err
	// }

	// defer res.Body.Close()
	// body, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	log.Printf("client: error reading response body: %s\n", err)
	// 	return err
	// }
	// var r model.DBResponse
	// err = json.Unmarshal(body, &r)
	// if err != nil {
	// 	log.Printf("client: error unmarshaling response body: %s\n", err)
	// 	return err
	// }
	// if r.Code != 200 {
	// 	log.Printf("client: error inserting message\n")
	// 	fmt.Println(r.Data, r.Code)
	// 	return fmt.Errorf("error inserting message: %s", r.Data)
	// }
	// log.Printf("Inserted %d message with ID %v", r.Data.InsertCount, r.Data.InsertId[0])
	// // log.Println(string(body))

	return nil
}

func (p *Profile) GetProfiles(chatID int64) ([]model.Profile, error) {
	conn := p.db

	expr := fmt.Sprintf("chat_id == %d", chatID)
	outputFields := []string{"profiles"}

	resq, err := conn.Query(
		context.Background(),    // context
		COLLECTION_NAME_PROFILE, // collectionName
		[]string{},              // partitionNames
		expr,                    // expr
		outputFields,            // outputFields
	)

	if err != nil {
		log.Println("Failed to get profiles: ", err.Error())
		return nil, err
	}

	return resultToJSONProfiles(resq), nil

}

func resultToJSONProfiles(results client.ResultSet) []model.Profile {
	var allProfiles []model.Profile

	// result only contains one colum with the profiles
	jsonData := results[0].FieldData().GetScalars().GetJsonData().Data

	var r model.ProfileResult
	if err := json.Unmarshal(jsonData[0], &r); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	allProfiles = append(allProfiles, r.Profiles...)

	return allProfiles

}

// define database schema for the "Profiles" collection
func getProfileSchema() *entity.Schema {
	chat_id := entity.NewField().
		WithName("chat_id").
		WithDataType(entity.FieldTypeInt64).
		WithIsPrimaryKey(true)

	vector := entity.NewField().
		WithName("vector").
		WithDataType(entity.FieldTypeFloatVector).
		WithDim(1)

	profiles := entity.NewField().
		WithName("profiles").
		WithDataType(entity.FieldTypeJSON)

	schema := &entity.Schema{
		CollectionName: COLLECTION_NAME_PROFILE,
		AutoID:         false,
		Fields: []*entity.Field{
			chat_id,
			vector,
			profiles,
		},
	}
	return schema
}
