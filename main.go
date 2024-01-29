package main

import (
	"example/plushie/plushie-bot/bot"
)

// var (
// 	bot *tgbotapi.BotAPI

// 	teddy      = "\U0001F9F8"
// 	r, _       = utf8.DecodeRuneInString(teddy)
// 	intro      = fmt.Sprintf("Hello! I am Plushie bot. %v You may talk to me about plushies!", string(r))
// 	statusText = "Plushies are happy and comfy."

// 	// Menu texts
// 	firstMenu  = fmt.Sprintf("<b>%v</b>\n\n%v", string(r), intro)
// 	secondMenu = "<b>How to use me</b>\n\nTry out any of the following commands:\n\n /return <i>plushie*n</i> - return n plushies\n /status - get status of plushies\n\n Or talk to me about plushies!"

// 	// Button texts
// 	startButton  = "Start"
// 	returnButton = "Return Plushie"
// 	statusButton = "Plushie Status"

// 	// tutorialButton = "Tutorial"

// 	// Keyboard layout for the first menu. One button, one row
// 	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData(startButton, startButton),
// 		),
// 	)

// 	// Keyboard layout for the second menu. Two buttons, one per row
// 	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData(returnButton, returnButton),
// 		),
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData(statusButton, statusButton),
// 		),
// 		// tgbotapi.NewInlineKeyboardRow(
// 		// 	tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
// 		// ),
// 	)
// )

func main() {

	/*
			db, err := database.NewPostgresDB(cfg.Database)
		if err != nil {
			logger.Fatal("failed connect to DB", zap.String("reason", err.Error()))
		}
		err = database.Migrate(&cfg.Migration, logger)
		if err != nil {
			logger.Fatal("can't run db migrations", zap.String("reason", err.Error()))
		}

			repo := repository.Init(db)
		svc := service.Init(repo)
		flows := flow.Init(svc)

		flows.ValidateCallbacksDataSize(&logger)

		bot := bot.Init(cfg, &logger, flows, svc, repo)

		bot.Run()
	*/
	bot := bot.Init()
	bot.Run()
}

/*
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
*/
