package repository

import (
	"context"
	"encoding/json"
	"example/bot/telegram-ai-bot/database"
	"example/bot/telegram-ai-bot/model"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

var COLLECTION_NAME_PROFILE = "Profiles"

type Profile struct {
	db client.Client
	// text string
}

func (p *Profile) Store(profiles *model.VectorizedProfile) error {
	jsonBytes, err := json.Marshal(profiles.Profiles)
	if err != nil {
		log.Printf("Error marshaling struct: %v", err)
	}
	profilesJSONObject := map[string]interface{}{
		"profiles": json.RawMessage(jsonBytes),
	}

	wrappedProfilesJSON, err := json.Marshal(profilesJSONObject)
	if err != nil {
		log.Fatal("Failed to marshal wrapped profiles:", err.Error())
	}

	// Prepare the row for insertion
	row := map[string]interface{}{
		"chat_id":  profiles.ChatID,
		"vector":   profiles.Vector,
		"profiles": wrappedProfilesJSON,
	}
	conn := database.Connect()

	// Insert the data
	_, err = conn.InsertRows(context.Background(), "Profiles", "", []interface{}{row})
	if err != nil {
		log.Fatal("Failed to insert rows:", err.Error())
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
