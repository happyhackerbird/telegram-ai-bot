package bot

import (
	"example/bot/telegram-ai-bot/controllers"
	"fmt"
	"log"
	"strings"

	"example/bot/telegram-ai-bot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) UpdateRouter(update tgbotapi.Update) {
	updLocal := model.DecodeToLocal(update)
	switch {
	case update.Message != nil:
		if strings.HasPrefix(update.Message.Text, "/") {
			b.SendMessage(b.CommandHandler(update.Message.Command(), updLocal))
		} else {
			b.SendMessage(b.MessageHandler(update.Message))
		}
		break
	case update.CallbackQuery != nil:
		b.SendMessage(b.CallbacksHandler(updLocal))
	}
}

// // We'll also say that this message is a reply to the previous message.
// // For any other specifications than Chat ID or Text, you'll need to
// // set fields on the `MessageConfig`.
// msg.ReplyToMessageID = update.Message.MessageID

func (b *Bot) MessageHandler(msg *tgbotapi.Message) tgbotapi.Chattable {
	chatID := msg.Chat.ID

	// handle profile creation
	if _, exists := b.userStates[chatID]; exists {
		return b.createProfile(msg.Text, chatID)
	}
	// handle AI response
	reply := tgbotapi.NewMessage(chatID, "")
	reply.Text = controllers.GetAIResponse(msg.Text)
	return reply
}

func (b *Bot) createProfile(userInput string, chatID int64) tgbotapi.Chattable {
	var msg tgbotapi.Chattable
	switch b.userStates[chatID] {
	case 0:
		b.UpdateProfile(chatID, "Name", userInput)
		b.userStates[chatID] = 1
		upd := model.UpdateLocal{TelegramChatID: model.TelegramChatID(chatID)}
		msg, _ = model.PromptInstructionHandler(&upd)
	case 1:
		b.UpdateProfile(chatID, "Instruction", userInput)
		b.userStates[chatID] = 2
		upd := model.UpdateLocal{TelegramChatID: model.TelegramChatID(chatID)}
		msg, _ = model.PromptAIModelHandler(&upd)
		// case "2":
		// 	b.UpdateProfile(chatID, "AIModel", userInput)
		//     finalizeProfileSetup(chatID, userState.ProfileData)
		//     // Remove user state as the process is complete
		// delete(userStates, chatID)
	}
	return msg
}

func (b *Bot) CommandHandler(cmd string, updLocal *model.UpdateLocal) tgbotapi.Chattable {
	chatId := int64(updLocal.TelegramChatID)
	var cd model.CallbackData
	switch cmd {
	case "start":
		cd = model.CallbackData{
			CommandKey: "start",
			Case:       "createProfile",
			Step:       0,
		}
	case "help":
		return b.sendMenu(chatId)
	case "profile":
		cd = model.CallbackData{CommandKey: "profile", Case: "options", Step: 0}
	default:
		return tgbotapi.NewMessage(chatId, "I don't know that command. Type /help for a help.")
	}

	// Handle the command using the flow system.
	replyMsg, err := b.Flow.Handle(&cd, updLocal)
	if err != nil {
		log.Printf("Error handling command: %s", err)
		return tgbotapi.NewMessage(chatId, "An error occurred.")
	}
	return replyMsg
}

func (b *Bot) CallbacksHandler(updLocal *model.UpdateLocal) tgbotapi.Chattable {
	// Decode the callback data from the query.
	cData := updLocal.CallbackData
	chatID := int64(updLocal.TelegramChatID)
	replyMessage, err := b.Flow.Handle(&cData, updLocal)
	if err != nil {
		log.Printf("Error handling callback query: %s", err)
		return tgbotapi.NewMessage(chatID, "An error occurred.")
	}

	// Acknowledge the callback query to prevent the loading spinner on the user's button.
	// callbackCfg := tgbotapi.NewCallback(, "")
	// if _, err := b.API.Send(callbackCfg); err != nil {
	// 	log.Printf("An error occurred acknowledging callback query: %s", err.Error())
	// }

	return replyMessage
}

func (b *Bot) StartProfileSetup(chatID int64) {
	b.userStates[chatID] = 0
}

func (b *Bot) UpdateProfile(chatId int64, key string, value string) {
	profile := b.Profiles[chatId]
	switch key {
	case "Name":
		profile.Name = value
	case "Instruction":
		profile.Instruction = value
	case "AIModel":
		profile.AIModel = value
	}
	b.Profiles[chatId] = profile
}

func (b *Bot) sendMenu(chatId int64) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, "")
	msg.Text = menu
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = menuMarkup
	return &msg
}

func (b *Bot) ShowProfile(msg *tgbotapi.MessageConfig, chatId int64) {
	profile := b.Profiles[chatId]
	msg.Text = fmt.Sprintf("<b>Profile</b>\n\nName: %v\nInstruction: %v\nAI Model: %v", profile.Name, profile.Instruction, profile.AIModel)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = menuMarkup
}
