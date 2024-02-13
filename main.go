package main

import (
	"example/bot/telegram-ai-bot/bot"
	"example/bot/telegram-ai-bot/flow"
	"example/bot/telegram-ai-bot/repository"
)

func main() {

	repo := repository.Init()
	flow := flow.Init()
	bot := bot.Init(flow, repo)
	bot.Run()
}
