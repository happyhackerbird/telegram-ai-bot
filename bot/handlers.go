package bot

import (
	"example/plushie/plushie-bot/controllers"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) UpdateRouter(update tgbotapi.Update) {
	if msg := update.Message; msg != nil {
		// Handle messages
		if update.Message != nil {
			if strings.HasPrefix(msg.Text, "/") {
				b.SendMessage(b.CommandHandler(msg.Command(), msg.Chat.ID))
			} else {
				b.SendMessage(b.ChatHandler(msg))
			}
		} else if update.CallbackQuery != nil {
			b.SendMessage(b.CallbacksHandler(update.CallbackQuery))
		}
	}

	// // We'll also say that this message is a reply to the previous message.
	// // For any other specifications than Chat ID or Text, you'll need to
	// // set fields on the `MessageConfig`.
	// msg.ReplyToMessageID = update.Message.MessageID

}

func (b *Bot) ChatHandler(msg *tgbotapi.Message) tgbotapi.Chattable {
	replyText := controllers.GetAIResponse(msg.Text)
	reply := tgbotapi.NewMessage(msg.Chat.ID, replyText)
	return reply
}

func (b *Bot) CommandHandler(cmd string, chatId int64) tgbotapi.Chattable {
	reply := tgbotapi.NewMessage(chatId, "")

	switch cmd {
	case "start":
		reply.Text = "Enter your custom AI assistant prompt."
	case "help":
		b.sendMenu(&reply, chatId)
	default:
		reply.Text = "I don't know that command. Type /help for help." + string(r)
	}
	return reply
}

func (b *Bot) CallbacksHandler(query *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	// cData := updLocal.CallbackData
	// replyMessage, err := b.Flow.Handle(&cData, updLocal)
	// if err != nil {
	// 	b.Logger.Error("error", zap.String("reason", err.Error()))
	// 	return commonanswers.UnknownError().BuildBotMessage(int64(updLocal.TelegramChatID))
	// }
	// return replyMessage

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, "")
	if query.Data == startButton {
		msg.Text = "Enter your custom AI assistant prompt."
	} else if query.Data == profileButton {
		msg.Text = profileText
	}
	callbackCfg := tgbotapi.NewCallback(query.ID, "")

	if _, err := b.API.Send(callbackCfg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}

	return msg
}

func (b *Bot) sendMenu(msg *tgbotapi.MessageConfig, chatId int64) {
	msg.Text = menu
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = menuMarkup
}
