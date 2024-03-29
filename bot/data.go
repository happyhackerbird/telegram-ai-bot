package bot

import (
	"example/bot/telegram-ai-bot/controllers"
	"example/bot/telegram-ai-bot/model"
	"example/bot/telegram-ai-bot/services"
	"fmt"
	"log"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	robot = "\U0001F916"
	r, _  = utf8.DecodeRuneInString(robot)

	intro       = fmt.Sprintf("Welcome to Build-A-Bot. %v You can create multiple custom AI chat profiles. Type /start to begin.", string(r))
	defaultText = "Custom AI profile not set. Create profile by typing /start."

	startButton   = "Start"
	profileButton = "Profile"

	menu       = "<b>Instructions</b>\n\n" + intro
	menuMarkup = tgbotapi.NewInlineKeyboardMarkup(
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

func (b *Bot) FinishProfileSetup(chatID int64) error {
	delete(b.userStates, chatID)

	//save profile in database
	p := b.Profiles[chatID]
	profiles := model.VectorizedProfile{
		ChatID:   chatID,
		Vector:   []float32{0},
		Profiles: []model.Profile{{Name: "Teddy Profile", Instruction: "Act as a small, playful and curious teddy bear. you have very soft feet and cuddly ears and squishy plush.", AIModel: "mixtral-8x7b-instruct"}, {Name: p.Name, Instruction: p.Instruction, AIModel: p.AIModel}},
	}
	delete(b.Profiles, chatID)

	err := b.Repository.Profile.Store(&profiles)
	if err != nil {
		log.Printf("Failed to insert message record: %s", err)
		return err
	}

	// user create profile, error on last step, retain profile ?
	return nil
}

func (b *Bot) SetModel(model string) {
	controllers.SetModel(model)
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
		controllers.SetInstruction(userInput)

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

	profiles, err := b.Repository.Profile.GetProfiles(chatId)
	if err != nil || profiles == nil {
		log.Printf("Failed to get profile record: %s", err)
		msg.Text = "An error occurred."
	}
	profile, exists := b.Profiles[chatId]
	if !exists && len(profiles) == 0 {
		msg.Text = defaultText
	} else {
		if exists {
			profiles = append(profiles, profile)
		}
		msg.Text = services.PrintProfiles(profiles)
	}
	msg.ParseMode = tgbotapi.ModeHTML
}

func (b *Bot) Store(msg *model.VectorizedMessage) {
	fmt.Println("Storing message in vector database ... ")
	err := b.Repository.Message.Store(msg)
	if err != nil {
		log.Printf("Failed to insert message record: %s", err)
	}
}
