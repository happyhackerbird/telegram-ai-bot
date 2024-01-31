package services

import (
	"example/bot/telegram-ai-bot/controllers"
	"example/bot/telegram-ai-bot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var b Bot

type Bot interface {
	UpdateProfile(chatID int64, field, value string)
	ShowProfile(msg *tgbotapi.MessageConfig, chatID int64)
	StartProfileSetup(chatID int64)
	FinishProfileSetup(chatID int64)
}

func SetBot(bot Bot) {
	b = bot
}

func PromptProfileNameHandler(updLocal *model.UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	b.StartProfileSetup(chatID)
	return tgbotapi.NewMessage(chatID, "Enter the name of the bot:"), nil
}

func PromptInstructionHandler(updLocal *model.UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	// txt := fmt.Sprintf("Set the name: %v. \n\n Enter the instruction", userInput)
	return tgbotapi.NewMessage(chatID, "Enter the instructions for the bot:"), nil
}

func PromptAIModelHandler(updLocal *model.UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	// txt := fmt.Sprintf("Set the instruction: %v. \n\n Select AI model", userInput)

	msg := tgbotapi.NewMessage(chatID, "Select the AI model:")
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", "start;createProfile;2;mixtral-8x7b-instruct"),
			tgbotapi.NewInlineKeyboardButtonData("2", "start;createProfile;2;pplx-70b-chat"),
			tgbotapi.NewInlineKeyboardButtonData("It's a surprise", "start;createProfile;2;codellama-70b-instruct"),
		),
	)
	return msg, nil
}

func FinalizeProfileHandler(updLocal *model.UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	b.UpdateProfile(chatID, "AIModel", updLocal.CallbackData.Payload)
	b.FinishProfileSetup(chatID)
	controllers.SetModel(updLocal.CallbackData.Payload) // where does this go

	return tgbotapi.NewMessage(int64(updLocal.TelegramChatID), "Profile created!"), nil
}

func ProfileOptionsHandler(updLocal *model.UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)
	msg := tgbotapi.NewMessage(chatID, "What do you want to do?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Create new profile", "start;createProfile;0;")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("View existing profiles", "profile;options;1;"),
		),
	)
	return msg, nil
}

func ViewProfilesHandler(updLocal *model.UpdateLocal) (tgbotapi.Chattable, error) {
	chatID := int64(updLocal.TelegramChatID)

	msg := tgbotapi.NewMessage(chatID, "")
	b.ShowProfile(&msg, chatID)
	return msg, nil
}
