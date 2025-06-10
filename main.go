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

	//c := cron.New()
	//
	//// Schedule to run daily at 9:00 AM
	//_, err = c.AddFunc("48 11 * * *", dailyBot.ProcessDailyChallenge)
	//if err != nil {
	//	logger.Fatalf("Failed to schedule job: %v", err)
	//}
	//
	//logger.Println("Starting LeetCode daily challenge bot")
	//c.Start()

	dailyBot.ProcessDailyChallenge()

	// Keep the application running
	select {}
}
