package config

// Config holds configuration data used by the application.
type Config struct {
	MinUsernameLength   int
	MinPassphraseLength int
	StorePath           string
	RequestTimeout      int
	SessionLength       int64
}

// NewConfiguration creates a new Config object with the default settings.
func NewConfiguration() Config {
	return Config{
		MinUsernameLength:   8,
		MinPassphraseLength: 16,
		StorePath:           "data/wasp.db",
		RequestTimeout:      30,      // 30 second time out
		SessionLength:       60 * 15, // 15 minute session
	}
}
