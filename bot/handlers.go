package bot

import (
	"encoding/json"
	"example/bot/telegram-ai-bot/controllers"
	"fmt"
	"log"
	"strings"

	"example/bot/telegram-ai-bot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) UpdateRouter(update tgbotapi.Update) {
	updLocal := model.DecodeToLocal(update)
	if msg := update.Message; msg != nil {
		// Handle messages
		if update.Message != nil {
			if strings.HasPrefix(msg.Text, "/") {
				b.SendMessage(b.CommandHandler(update.Message.Command(), updLocal))
			} else {
				b.SendMessage(b.MessageHandler(updLocal))
			}
		}
		if update.CallbackQuery != nil {
			b.SendMessage(b.CallbacksHandler(update.CallbackQuery))
		}
	}

	// // We'll also say that this message is a reply to the previous message.
	// // For any other specifications than Chat ID or Text, you'll need to
	// // set fields on the `MessageConfig`.
	// msg.ReplyToMessageID = update.Message.MessageID

}

func (b *Bot) MessageHandler(upd tgbotapi.Update) tgbotapi.Chattable {
	updLocal := model.DecodeToLocal(upd)

	replyText := controllers.GetAIResponse(msg.Text)
	reply := tgbotapi.NewMessage(msg.Chat.ID, replyText)
	return reply
}

func (b *Bot) CommandHandler(cmd string, updLocal *model.UpdateLocal) tgbotapi.Chattable {
	var cd model.CallbackData
	switch cmd {
	case "start":
		cd = model.CallbackData{
			CommandKey: "start",
			Case:       "createProfile",
			Step:       0,
		}
	case "help":
		return b.sendMenu(updLocal.TelegramChatID)
	case "profile":
		cd = model.CallbackData{CommandKey: "profile", Case: "options", Step: 0}
	default:
		return tgbotapi.NewMessage(updLocal.TelegramChatID, "I don't know that command. Type /help for a help.")
	}
	// Handle the command using the flow system.
	replyMsg, err := model.Flow.Handle(&cd, updLocal)
	if err != nil {
		log.Printf("Error handling command: %s", err)
		return tgbotapi.NewMessage(updLocal.TelegramChatID, "An error occurred.")
	}
	return replyMsg
}

func (b *Bot) CallbacksHandler(updLocal *model.UpdateLocal) tgbotapi.Chattable {
	// Decode the callback data from the query.
	var cd model.CallbackData
	if err := json.Unmarshal([]byte(query.Data), &cd); err != nil {
		log.Printf("Error decoding callback data: %s", err)
		return tgbotapi.NewMessage(query.Message.Chat.ID, "An error occurred.")
	}

	// Create an UpdateLocal object, assuming it contains necessary information like chat ID.
	updLocal := &model.UpdateLocal{TelegramChatID: query.Message.Chat.ID}

	// Handle the callback query using the flow system.
	replyMsg, err := model.Flow.Handle(&cd, updLocal)
	if err != nil {
		log.Printf("Error handling callback query: %s", err)
		return tgbotapi.NewMessage(query.Message.Chat.ID, "An error occurred.")
	}

	// Acknowledge the callback query to prevent the loading spinner on the user's button.
	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	if _, err := b.API.Send(callbackCfg); err != nil {
		log.Printf("An error occurred acknowledging callback query: %s", err.Error())
	}

	return replyMsg
}

func (b *Bot) sendMenu(chatId int64) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, "")
	msg.Text = menu
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = menuMarkup
	return &msg
}

func (b *Bot) showProfile(msg *tgbotapi.MessageConfig, chatId int64) {
	profile := b.Profiles[chatId]
	msg.Text = fmt.Sprintf("<b>Profile</b>\n\nName: %v\nInstruction: %v\nAI Model: %v", profile.Name, profile.Instruction, profile.AIModel)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = menuMarkup
}
