package config

// Config holds configuration data used by the application.
type Config struct {
	MinUsernameLength   int
	MinPassphraseLength int
	ArgonTime           uint32
	ArgonMemory         uint32
	ArgonThreads        uint8
	StorePath           string
	RequestTimeout      int
	SessionLength       int64
	TotpStart           int64
	TotpStep            int64
	TotpLength          int
	TotpName            string
	TotpIssuer          string
}

// NewConfiguration creates a new Config object with the default settings.
func NewConfiguration() Config {
	return Config{
		MinUsernameLength:   8,
		MinPassphraseLength: 16,
		ArgonTime:           4,
		ArgonMemory:         64 * 1024,
		ArgonThreads:        3,
		StorePath:           "data/wasp.db",
		RequestTimeout:      30,      // 30 second time out
		SessionLength:       60 * 15, // 15 minute session
		TotpStart:           0,       // Start at the Unix epoch
		TotpStep:            30,      // Get a new code every 30 seconds
		TotpLength:          6,       // Totp value should be 6 digits
		TotpName:            "wasp",
		TotpName:            "wasp",
	}
}
