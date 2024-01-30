package main

import (
	"example/bot/telegram-ai-bot/bot"
)

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
