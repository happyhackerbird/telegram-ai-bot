package bot

import (
	"example/bot/telegram-ai-bot/controllers"
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
	case update.CallbackQuery != nil:
		b.SendMessage(b.CallbacksHandler(updLocal, update.CallbackQuery.ID))
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
	reply.Text = controllers.GetAIResponse(chatID, msg.Text)
	return reply
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
		return tgbotapi.NewMessage(chatId, "I don't know that command. Type /help for help.")
	}

	// Handle the command using the flow system.
	replyMsg, err := b.Flow.Handle(&cd, updLocal)
	if err != nil {
		log.Printf("Error handling command: %s", err)
		return tgbotapi.NewMessage(chatId, "An error occurred.")
	}
	return replyMsg
}

func (b *Bot) CallbacksHandler(updLocal *model.UpdateLocal, id string) tgbotapi.Chattable {
	// Decode the callback data from the query.
	cData := updLocal.CallbackData
	chatID := int64(updLocal.TelegramChatID)
	replyMessage, err := b.Flow.Handle(&cData, updLocal)
	if err != nil {
		log.Printf("Error handling callback query: %s", err)
		return tgbotapi.NewMessage(chatID, "An error occurred.")
	}

	// Acknowledge the callback query to prevent the loading spinner on the user's button.
	callbackCfg := tgbotapi.NewCallback(id, "")
	if _, err := b.API.Send(callbackCfg); err != nil {
		log.Printf("An error occurred acknowledging callback query: %s", err.Error())
	}

	return replyMessage
}
