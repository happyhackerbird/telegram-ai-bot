package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"example/plushie/plushie-bot/controllers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	bot *tgbotapi.BotAPI

	teddy      = "\U0001F9F8"
	r, _       = utf8.DecodeRuneInString(teddy)
	intro      = fmt.Sprintf("Hello! I am Plushie bot. %v You may talk to me about plushies!", string(r))
	statusText = "Plushies are happy and comfy."

	// Menu texts
	firstMenu  = fmt.Sprintf("<b>%v</b>\n\n%v", string(r), intro)
	secondMenu = "<b>How to use me</b>\n\nTry out any of the following commands:\n\n /return <i>plushie*n</i> - return n plushies\n /status - get status of plushies\n\n Or talk to me about plushies!"

	// Button texts
	startButton  = "Start"
	returnButton = "Return Plushie"
	statusButton = "Plushie Status"

	// tutorialButton = "Tutorial"

	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(startButton, startButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(returnButton, returnButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(statusButton, statusButton),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
		// ),
	)
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	// Set the bot to use debug mode (verbose logging).
	bot.Debug = false
	log.Printf("Authorized as @%s. %v", bot.Self.UserName, string(r))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

	// // Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// // to make sure Telegram knows we've handled previous values and we don't
	// // need them repeated.
	// updateConfig := tgbotapi.NewUpdate(0)

	// // Tell Telegram we should wait up to 30 seconds on each request for an
	// // update. This way we can get information just as quickly as making many
	// // frequent requests without having to send nearly as many.
	// updateConfig.Timeout = 30

	// // Start polling Telegram for updates.
	// updates := bot.GetUpdatesChan(updateConfig)

	// // // Let's go through each update that we're getting from Telegram.
	// for update := range updates {
	// 	// Telegram can send many types of updates depending on what your Bot
	// 	// is up to. We only want to look at messages for now, so we can
	// 	// discard	 any other updates.
	// 	if update.Message == nil {
	// 		continue
	// 	}

	// }
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break
	}
}

func handleMessage(msg *tgbotapi.Message) {
	var err error
	if strings.HasPrefix(msg.Text, "/") {
		err = handleCommand(msg)
	} else {
		err = handleChatQuery(msg)
	}
	// // We'll also say that this message is a reply to the previous message.
	// // For any other specifications than Chat ID or Text, you'll need to
	// // set fields on the `MessageConfig`.
	// msg.ReplyToMessageID = update.Message.MessageID

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func handleChatQuery(msg *tgbotapi.Message) error {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "")
	reply.Text = controllers.GetAIResponse(msg.Text)
	_, err := bot.Send(reply)
	return err
}

func handleCommand(msg *tgbotapi.Message) error {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "")
	// Extract the command from the Message.
	switch msg.Command() {
	case "start":
		reply.Text = intro + "\n\n" + "Type /help to see what I can do."
	case "status":
		reply.Text = statusText
	case "return":
		controllers.ReturnPlushies(msg.Text, &reply)
	case "help":
		sendMenu(&reply, msg.Chat.ID)
	default:
		reply.Text = "I don't know that command. Here is a plushie: " + string(r)
	}
	_, err := bot.Send(reply)
	return err
}

func handleButton(query *tgbotapi.CallbackQuery) {
	var msg tgbotapi.MessageConfig

	if query.Data == startButton {
		msg := tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, secondMenu, secondMenuMarkup)
		msg.ParseMode = tgbotapi.ModeHTML
		if _, err := bot.Send(msg); err != nil {
			log.Printf("An error occured: %s", err.Error())
		}
		return
	}

	msg = tgbotapi.NewMessage(query.Message.Chat.ID, "")
	if query.Data == returnButton {
		controllers.ReturnRandomPlushies(&msg)
	} else if query.Data == statusButton {
		msg.Text = statusText
	}
	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	bot.Send(callbackCfg)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}

}

func sendMenu(msg *tgbotapi.MessageConfig, chatId int64) {
	msg.Text = firstMenu
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
}
