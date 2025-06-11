package main

import (
	"log"
	"os"

	"github.com/thuanvo2008/leet-code-daily-bot/service"
)

func main() {
	logger := log.New(os.Stdout, "[DailyBot] ", log.LstdFlags)

	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	dailyBot, err := service.NewService(webhookURL, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize service: %v", err)
	}

	dailyBot.ProcessDailyChallenge()
}
