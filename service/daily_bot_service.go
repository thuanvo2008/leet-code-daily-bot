package service

import (
	"github.com/thuanvo2008/leet-code-daily-bot/client"
	error2 "github.com/thuanvo2008/leet-code-daily-bot/error"

	"log"
	"net/http"
)

// Service manages the application workflow
type Service struct {
	leetcodeClient *client.LeetCodeClient
	discordClient  *client.DiscordClient
	webhookURL     string
	logger         *log.Logger
}

// NewService creates a new application service
func NewService(webhookURL string, logger *log.Logger) (*Service, error) {
	if webhookURL == "" {
		return nil, error2.ErrMissingWebhookURL
	}

	httpClient := &http.Client{}
	return &Service{
		leetcodeClient: client.NewLeetCodeClient(httpClient),
		discordClient:  client.NewDiscordClient(httpClient, webhookURL),
		logger:         logger,
	}, nil
}

// ProcessDailyChallenge fetches and posts the daily challenge
func (s *Service) ProcessDailyChallenge() {
	s.logger.Println("Processing daily challenge")

	resp, err := s.leetcodeClient.FetchDailyProblem()
	if err != nil {
		s.logger.Printf("Error fetching problem: %v", err)
		return
	}

	question := &resp.Data.ActiveDailyCodingChallengeQuestion.Question
	date := resp.Data.ActiveDailyCodingChallengeQuestion.Date

	err = s.discordClient.PostDailyChallenge(question, date)

	if err != nil {
		s.logger.Printf("Error posting to Discord: %v", err)
		return
	}

	s.logger.Println("Successfully processed daily challenge")
}
