package main

import (
	"example/bot/telegram-ai-bot/bot"
	"example/bot/telegram-ai-bot/database"
	"example/bot/telegram-ai-bot/flow"
	"example/bot/telegram-ai-bot/repository"
)

func main() {
	db := database.Connect()
	database.LoadIndex()
	repo := repository.Init(db)
	flow := flow.Init()
	bot := bot.Init(flow, repo)
	bot.Run()
}
