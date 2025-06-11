package client

import (
	"encoding/json"
	"fmt"
	"github.com/thuanvo2008/leet-code-daily-bot/model"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// DiscordClient handles communication with Discord webhooks
type DiscordClient struct {
	requestSender RequestSender
	webhookURL    string
}

// NewDiscordClient creates a new client for Discord interactions
func NewDiscordClient(client *http.Client, webhookUrl string) *DiscordClient {
	return &DiscordClient{
		requestSender: NewHTTPRequestClient(client),
		webhookURL:    webhookUrl,
	}
}

// PostDailyChallenge sends a daily challenge update to Discord
func (c *DiscordClient) PostDailyChallenge(question *model.Question, date string) error {
	payload, err := c.createDiscordPayload(question, date)
	if err != nil {
		return fmt.Errorf("failed to create payload: %w", err)
	}

	req, err := c.createWebhookRequest(payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	_, err = c.requestSender.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to post to Discord: %w", err)
	}

	// Success case
	fmt.Printf("Successfully posted: %s (%s)\n", question.Title, question.Difficulty)
	return nil
}

// createDiscordPayload builds a Discord message payload
func (c *DiscordClient) createDiscordPayload(question *model.Question, date string) ([]byte, error) {
	// Format the date if needed
	formattedDate := fmt.Sprintf("📌 # LeetCode Daily Challenge (%s)", formatDate(date))
	payload := model.DiscordPayload{
		Content: fmt.Sprintf("@everyone 🎅Hey noob: \n\n %s \n\n %s", formattedDate, c.formatDescription(question)),
	}

	return json.Marshal(payload)
}

func (c *DiscordClient) formatDescription(question *model.Question) string {
	diffEmoji := "🧮"
	switch question.Difficulty {
	case "Easy":
		diffEmoji = "🟢"
	case "Medium":
		diffEmoji = "🟠"
	case "Hard":
		diffEmoji = "🔴"
	}

	// Format the content using our HTML formatter
	formattedContent := c.formatHTMLContent(question.Content)

	// Format the description nicely
	return fmt.Sprintf("📺 %s %s (%s)\n🔗 %s\n\n📝 Question Detail:\n-----------\n👉 %s",
		diffEmoji,
		question.Title,
		question.Difficulty,
		question.Url,
		formattedContent)
}

// getDifficultyColor returns a color code based on difficulty
func (c *DiscordClient) getDifficultyColor(difficulty string) int {
	switch difficulty {
	case "Easy":
		return 5025616 // Green
	case "Medium":
		return 16750848 // Orange
	case "Hard":
		return 16007990 // Red
	default:
		return 5793266 // Blue
	}
}

// createWebhookRequest creates an HTTP request for Discord webhook
func (c *DiscordClient) createWebhookRequest(payload []byte) (*http.Request, error) {
	req, err := CreateRequest("POST", c.webhookURL, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// createThread creates a thread on a message
func (c *DiscordClient) createThread(messageID, title string) error {
	// Extract channel ID from webhook URL
	parts := strings.Split(c.webhookURL, "/")
	channelID := parts[len(parts)-2]

	// Create thread payload
	threadPayload := model.ThreadCreatePayload{
		Name:                fmt.Sprintf("Discussion: %s", title),
		AutoArchiveDuration: 1440, // 1 day in minutes
	}

	payloadBytes, err := json.Marshal(threadPayload)
	if err != nil {
		return err
	}

	// Discord API endpoint for creating a thread from a message
	apiURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s/threads",
		channelID, messageID)

	req, err := CreateRequest("POST", apiURL, payloadBytes)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", os.Getenv("DISCORD_BOT_TOKEN")))

	_, err = c.requestSender.SendRequest(req)
	return err
}

// formatDate ensures the date is in YYYY-MM-DD format
func formatDate(date string) string {
	// Parse the date if needed
	parsedDate, err := time.Parse("2006-01-02", date)
	if err == nil {
		return parsedDate.Format("2006-01-02")
	}
	return date
}

// Add this function to your discord_client.go file
func (c *DiscordClient) formatHTMLContent(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}

	// Replace common HTML elements with Discord markdown
	content := htmlContent

	// Replace paragraph tags with newlines
	content = strings.ReplaceAll(content, "<p>", "")
	content = strings.ReplaceAll(content, "</p>", "\n")

	// Replace code blocks
	content = strings.ReplaceAll(content, "<pre>", "```\n")
	content = strings.ReplaceAll(content, "</pre>", "\n```")
	content = strings.ReplaceAll(content, "<code>", "`")
	content = strings.ReplaceAll(content, "</code>", "`")

	// Replace lists
	content = strings.ReplaceAll(content, "<ul>", "")
	content = strings.ReplaceAll(content, "</ul>", "")
	content = strings.ReplaceAll(content, "<li>", "• ")
	content = strings.ReplaceAll(content, "</li>", "\n")

	// Replace bold and italic
	content = strings.ReplaceAll(content, "<strong>", "**")
	content = strings.ReplaceAll(content, "<strong class=\"example\">", "👉**")
	content = strings.ReplaceAll(content, "</strong>", "**")
	content = strings.ReplaceAll(content, "<em>", "")
	content = strings.ReplaceAll(content, "</em>", "")
	content = strings.ReplaceAll(content, "<sub>", "")
	content = strings.ReplaceAll(content, "</sub>", "")
	content = strings.ReplaceAll(content, "<span class=\"example-io\">", "")
	content = strings.ReplaceAll(content, "</span>", "")
	content = strings.ReplaceAll(content, "<div class=\"example-block\">", "")
	content = strings.ReplaceAll(content, "</div>", "")

	// Replace common HTML entities
	content = strings.ReplaceAll(content, "&lt;", "<")
	content = strings.ReplaceAll(content, "&gt;", ">")
	content = strings.ReplaceAll(content, "&amp;", "&")
	content = strings.ReplaceAll(content, "&quot;", "\"")
	content = strings.ReplaceAll(content, "&nbsp;", "-----------------------------=================================----------------------------------")
	content = strings.ReplaceAll(content, "&#39;", "'")
	content = strings.ReplaceAll(content, "<font face=\"monospace\">", "")
	content = strings.ReplaceAll(content, "</font>", "")
	content = strings.ReplaceAll(content, "</div>", "")

	// Clean up excessive newlines
	content = strings.ReplaceAll(content, "\n\n\n", "\n\n")

	// Remove any remaining HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	content = re.ReplaceAllString(content, "")

	// Handle multiple empty lines more effectively
	// First convert 3+ consecutive newlines to just 2 newlines (one blank line)
	multipleNewlines := regexp.MustCompile(`\n{3,}`)
	content = multipleNewlines.ReplaceAllString(content, "\n\n")

	// Then fix leading spaces at beginning of non-empty lines
	leadingSpaces := regexp.MustCompile(`(?m)^[ \t]+`)
	content = leadingSpaces.ReplaceAllString(content, "")

	// Remove whitespace-only lines (lines with just spaces/tabs)
	emptyLines := regexp.MustCompile(`(?m)^[ \t]*\n`)
	content = emptyLines.ReplaceAllString(content, "\n")

	// Finally ensure we don't have more than 2 consecutive newlines anywhere
	content = multipleNewlines.ReplaceAllString(content, "\n\n")

	// Discord has a limit of 4096 characters per embed description
	if len(content) > 4000 {
		return content[:3997] + "..."
	}

	return content
}
