package model

import (
	// "encoding/json"
	// "fmt"
	// "log"

	"fmt"

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

type Bot interface {
	UpdateProfile(chatID int64, field, value string)
	ShowProfile(msg *tgbotapi.MessageConfig, chatID int64)
	StartProfileSetup(chatID int64)
}

var b Bot

func SetBot(bot Bot) {
	b = bot
}

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

// here is just a way to get a handler at the service level
func (flow Flow) Handle(cd *CallbackData, updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	msg, err := flow[cd.CommandKey][cd.Case][cd.Step].Handler(updLocal)
	fmt.Println("handling ", flow[cd.CommandKey][cd.Case][cd.Step])
	if err != nil {
		return nil, err
	}
	// if cd.CommandKey == "start" {
	// 	chatID := int64(updLocal.TelegramChatID)

	// // Store the user input after each step
	// switch cd.Step {
	// case 0:
	//    ProfileNameResponseHandler(updLocal, updLocal.)
	// case 1:
	// 	b.UpdateProfile(chatID, "Instruction", updLocal.CallbackData.Payload)
	// case 2:
	// 	b.UpdateProfile(chatID, "AIModel", updLocal.CallbackData.Payload)
	// }
	// } else if cd.CommandKey == "profile" {
	// 	switch cd.Case {
	// 	case "options":
	// 		switch cd.Step {
	// 		case 2:
	// 			// Handle the viewing of existing profiles
	// 			b.showProfile(msg, updLocal.TelegramChatID)
	// 		}
	// 	}
	// }

	return msg, nil
}

func PromptProfileNameHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	b.StartProfileSetup(chatID)
	return tgbotapi.NewMessage(chatID, "Enter the name of the bot"), nil
}

func PromptInstructionHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	// txt := fmt.Sprintf("Set the name: %v. \n\n Enter the instruction", userInput)
	return tgbotapi.NewMessage(chatID, "Enter the instructions for the bot"), nil
}

func PromptAIModelHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	// txt := fmt.Sprintf("Set the instruction: %v. \n\n Select AI model", userInput)

	msg := tgbotapi.NewMessage(chatID, "Select the AI model")
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Normal", "start;createProfile;2;mixtral-8x7b-instruct"),
			tgbotapi.NewInlineKeyboardButtonData("Creative", "start;createProfile;2;llama-2-70b-chat"),
		),
	)
	return msg, nil
}

func StoreAIModelHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	b.UpdateProfile(int64(updLocal.TelegramChatID), "AIModel", updLocal.CallbackData.Payload)
	return tgbotapi.NewMessage(int64(updLocal.TelegramChatID), "Profile created"), nil
}

func ProfileOptionsHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	msg := tgbotapi.NewMessage(chatID, "What do you want to do?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Create new profile", "createProfile")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("View existing profiles", "viewProfiles"),
		),
	)
	return msg, nil
}

func ViewProfilesHandler(updLocal *UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)

	msg := tgbotapi.NewMessage(chatID, "")
	b.ShowProfile(&msg, chatID)
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
