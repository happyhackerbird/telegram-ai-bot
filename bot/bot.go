package bot

import (
	"bufio"
	"context"
	"example/plushie/plushie-bot/model"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

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
	bot   Bot
	robot = "\U0001F916"
	r, _  = utf8.DecodeRuneInString(robot)

	intro       = fmt.Sprintf("Hello! I am your AI assistant. %v You can configure me with a custom prompt. Type /start to begin.", string(r))
	profileText = "Custom AI profile not set. Enter your custom AI assistant prompt by typing start."

	startButton   = "Start"
	profileButton = "Profile"
	menu          = "<b>Instructions for the human</b>\n\n" + intro
	menuMarkup    = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(startButton, startButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(profileButton, profileButton),
		),
	)
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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go b.receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

func (b *Bot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			b.UpdateRouter(update)
		}
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

func (b *Bot) SendMessage(msg tgbotapi.Chattable) {
	_, err := b.API.Request(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
