package model

// DiscordEmbed represents a Discord message embed
type DiscordEmbed struct {
	Title       string        `json:"title"`
	URL         string        `json:"url"`
	Description string        `json:"description"`
	Timestamp   string        `json:"timestamp"`
	Footer      DiscordFooter `json:"footer"`
	Color       int           `json:"color"` // Color in decimal format
}

type DiscordFooter struct {
	Text string `json:"text"`
}

// DiscordPayload represents a Discord webhook payload
type DiscordPayload struct {
	Embeds  []DiscordEmbed `json:"embeds"`
	Content string         `json:"content"`
}

// DiscordResponse represents the response from Discord webhook
type DiscordResponse struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
}

// ThreadCreatePayload represents payload for thread creation
type ThreadCreatePayload struct {
	Name                string `json:"name"`
	AutoArchiveDuration int    `json:"auto_archive_duration"`
}
