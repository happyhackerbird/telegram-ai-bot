package bot

import (
	"example/bot/telegram-ai-bot/database"
	"example/bot/telegram-ai-bot/model"
	"example/bot/telegram-ai-bot/services"
	"fmt"
	"log"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Profile struct {
	Name        string
	Instruction string
	AIModel     string
}

var (
	robot = "\U0001F916"
	r, _  = utf8.DecodeRuneInString(robot)

	intro       = fmt.Sprintf("Hello! I am your AI assistant. %v You can configure me with a custom prompt. Type /start to begin.", string(r))
	defaultText = "Custom AI profile not set. Create profile by typing /start."

	startButton   = "Start"
	profileButton = "Profile"
	menu          = "<b>Instructions for the human</b>\n\n" + intro
	menuMarkup    = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(startButton, "start;createProfile;0;"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(profileButton, "profile;options;1;"),
		),
	)
)

func (b *Bot) sendMenu(chatId int64) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, "")
	msg.Text = menu
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = menuMarkup
	return &msg
}

func (b *Bot) StartProfileSetup(chatID int64) {
	b.userStates[chatID] = 0
}

func (b *Bot) FinishProfileSetup(chatID int64) {
	delete(b.userStates, chatID)
}

func (b *Bot) createProfile(userInput string, chatID int64) tgbotapi.Chattable {
	var msg tgbotapi.Chattable
	switch b.userStates[chatID] {
	case 0:
		b.UpdateProfile(chatID, "Name", userInput)
		b.userStates[chatID] = 1
		upd := model.UpdateLocal{TelegramChatID: model.TelegramChatID(chatID)}
		msg, _ = services.PromptInstructionHandler(&upd)
	case 1:
		b.UpdateProfile(chatID, "Instruction", userInput)
		b.userStates[chatID] = 2
		upd := model.UpdateLocal{TelegramChatID: model.TelegramChatID(chatID)}
		database.SetInstruction(userInput) //unclean

		msg, _ = services.PromptAIModelHandler(&upd)
	}
	return msg
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

func (b *Bot) ShowProfile(msg *tgbotapi.MessageConfig, chatId int64) {
	if profile, exists := b.Profiles[chatId]; !exists {
		msg.Text = defaultText
	} else {
		msg.Text = fmt.Sprintf("<b>Profile</b>\n\nName: %v\nInstruction: %v\nAI Model: %v", profile.Name, profile.Instruction, profile.AIModel)
	}
	msg.ParseMode = tgbotapi.ModeHTML
}

func (b *Bot) GetCount() int64 {
	c := b.message_count
	b.message_count++
	return c
}

func (b *Bot) DiscardCount() {
	b.message_count--
}

func (b *Bot) Store(msg *model.VectorizedMessage) {
	fmt.Println("Storing message in vector database ... ")
	err := b.Repository.Message.Store(msg)
	if err != nil {
		log.Printf("Failed to insert message record: %s", err)
		b.DiscardCount()
	}
}
