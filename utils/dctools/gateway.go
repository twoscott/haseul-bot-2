package dctools

// NewBot prefixes the Bot token type to the provided token.
func BotToken(token string) string {
	return "Bot " + token
}
