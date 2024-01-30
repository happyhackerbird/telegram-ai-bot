package model

import (
	"encoding/json"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
CommandFlow (starting a specific script (flow) for manipulating an object, is a command)
		  |
		  Usecase (actions that can be performed on an object)
				|
				Chain (algorithm, sequence of steps to implement an action and obtain some result)
					|
					Step (certain, specific step, action)
*/

type (
	Flow    map[CommandKey]Usecase
	Usecase map[Case]Chain
	Case    string
	Chain   map[Step]Action
	Step    int
	Action  struct {
		Handler HandlerFunc
		Message
	}
	HandlerFunc func(*UpdateLocal) (tgbotapi.Chattable, error)
	Message     struct {
		Text    string
		Buttons []Button
	}
	Button struct {
		Name         string
		CallbackData CallbackData
	}
)

var flow Flow
var jsonData = `
{
   "start":{
      "createProfile":{
         "0":{
            "handler": PromptProfileNameHandler(),
            "message":"Enter the name of the bot",
            "buttons": nil
         },
         "1":{
            "handler": PromptInstructionHandler(),
            "message":"Enter the instruction",
            "buttons": nil
         },
         "2":{
            "handler": PromptAIModelHandler(),
            "message":"Select AI model",
            "buttons": [
               {
                  "name":"Normal",
                  "callback_data":{
                     "cmd_key":"start",
                     "case":"createProfile",
                     "step":2,
                     "payload":"mixtral-8x7b-instruct"
                  }
               },
               {
                  "name":"Creative",
                  "callback_data":{
                     "cmd_key":"start",
                     "case":"createProfile",
                     "step":2,
                     "payload":"llama-2-70b-chat"
                  }
               }
            ]
         }
      }
   },
   "profile":{
      "options":{
         "0":{
            "handler": ProfileOptionsHandler(),
            "message":"What do you want to do?",
            "buttons":[
               {
                  "name":"Create new profile",
                  "callback_data":{
                     "cmd_key":"start",
                     "case":"createProfile",
                     "step":0,
                     "payload":""
                  }
               },
               {
                  "name":"View existing profiles",
                  "callback_data":{
                     "cmd_key":"profile",
                     "case":"options",
                     "step":2,
                     "payload":"view"
                  }
               }
            ]
         },
         "2":{
            "handler": ViewProfilesHandler(),
            "message":"", // The message will be set in the handler
            "buttons": nil
         }
      }
   }
}`

func init() {
	err := json.Unmarshal([]byte(jsonData), &flow)
	if err != nil {
		log.Fatalf("error loading flow data: %v", err)
	}
}

// here is just a way to get a handler at the service level
func (flow Flow) Handle(cd *CallbackData, updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	msg, err := flow[cd.CommandKey][cd.Case][cd.Step].Handler(updLocal)
	if err != nil {
		return nil, err
	}
	if cd.CommandKey == "start" {
		// Store the user input after each step
		switch cd.Step {
		case 0:
			b.updateProfile(updLocal, "Name", updLocal.Message.Text)
		case 1:
			b.updateProfile(updLocal, "Instruction", updLocal.Message.Text)
		case 2:
			b.updateProfile(updLocal, "AIModel", updLocal.Message.Text)
		}
		// } else if cd.CommandKey == "profile" {
		// 	switch cd.Case {
		// 	case "options":
		// 		switch cd.Step {
		// 		case 2:
		// 			// Handle the viewing of existing profiles
		// 			b.showProfile(msg, updLocal.TelegramChatID)
		// 		}
		// 	}
	}

	return msg, nil
}

func PromptProfileNameHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(updLocal.TelegramChatID, "Enter the name of the bot"), nil
}

func PromptInstructionHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(updLocal.TelegramChatID, "Enter the instruction"), nil
}

func PromptAIModelHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(updLocal.TelegramChatID, "Select AI model"), nil
}

func ProfileOptionsHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(updLocal.TelegramChatID, "What do you want to do?"), nil
}

func ViewProfilesHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(updLocal.TelegramChatID, "")
	b.showProfile(&msg, updLocal.TelegramChatID)
	return msg, nil
}

// // a general function for assembling a bot message from the described local model
// func (msg Message) BuildBotMessage(chatID int64) tgbotapi.MessageConfig {
// 	replyMessage := tgbotapi.NewMessage(chatID, msg.Text)
// 	var buttonRows [][]tgbotapi.InlineKeyboardButton
// 	for _, button := range msg.Buttons {
// 		buttonRows = append(buttonRows, tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData(button.Name, button.CallbackData.Encode()),
// 		),
// 		)
// 	}
// 	markup := tgbotapi.NewInlineKeyboardMarkup(
// 		buttonRows...,
// 	)
// 	replyMessage.ReplyMarkup = markup
// 	replyMessage.ParseMode = tgbotapi.ModeHTML
// 	return replyMessage
// }

// // helper function for checking the size of the described callback data at the start of the application
// func (flow *Flow) ValidateCallbacksDataSize(logger *zap.Logger) {
// 	for _, usecase := range *flow {
// 		for _, chain := range usecase {
// 			for _, action := range chain {
// 				for _, button := range action.Buttons {
// 					// 64 bytes - telegram limit for callback_data: https://core.telegram.org/bots/api#inlinekeyboardbutton
// 					if len(button.CallbackData.Encode()) > 64 {
// 						logger.Fatal("size of callback_data exceeds 64 bytes", zap.String("callback_data", button.CallbackData.Encode()), zap.Int("bytes", len(button.CallbackData.Encode())))
// 					}
// 				}
// 			}
// 		}
// 	}
// 	logger.Info("callback_data dimensions are valid")
// }

// Helper function to update the profile
func (b *Bot) updateProfile(updLocal *UpdateLocal, key string, value string) {
	profile := b.Profiles[updLocal.TelegramChatID]
	switch key {
	case "Name":
		profile.Name = value
	case "Instruction":
		profile.Instruction = value
	case "AIModel":
		profile.AIModel = value
	}
	b.Profiles[updLocal.TelegramChatID] = profile
}
