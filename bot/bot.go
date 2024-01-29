package bot

import (
	"example/plushie/plushie-bot/model"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Bot struct {
	API *tgbotapi.BotAPI
	// Config     *config.Config
	// Logger     *zap.Logger
	// Flow model.Flow
	// Service    *service.Service
	// Repository *repository.Repository
}

var (
	bot Bot
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Init() Bot {
	return Bot{
		API: nil,
	}
}

func (b *Bot) Run() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	log.Printf("Authorized as @%s.", bot.Self.UserName)
	b.API = bot

	// Set the bot to use debug mode (verbose logging).
	bot.Debug = false

	err = b.SetBotCommands()
	if err != nil {
		log.Fatalf("Error setting bot commands: %v", err)
	}

}

// configure the bot menu, don't use "start" command, but you can if you want
func (b *Bot) InitBotCommands() tgbotapi.SetMyCommandsConfig {
	commands := []model.CommandEntity{
		{
			Key:  model.ProfileCommand,
			Name: "profile",
		},
		{
			Key:  model.StartCommand,
			Name: "start",
		},
	}
	tgCommands := make([]tgbotapi.BotCommand, 0, len(commands))
	for _, cmd := range commands {
		tgCommands = append(tgCommands, tgbotapi.BotCommand{
			Command:     "/" + string(cmd.Key),
			Description: cmd.Name,
		})
	}
	commandsConfig := tgbotapi.NewSetMyCommands(tgCommands...)
	return commandsConfig
}

func (b *Bot) SetBotCommands() error {
	commandsConfig := b.InitBotCommands()
	_, err := b.API.Request(commandsConfig)
	if err != nil {
		return err
	}
	return nil
}
