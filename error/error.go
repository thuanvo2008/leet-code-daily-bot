package error

// Custom errors
var ErrMissingWebhookURL = Error("DISCORD_WEBHOOK_URL is not set")

// Error represents an application-specific error
type Error string

func (e Error) Error() string {
	return string(e)
}
