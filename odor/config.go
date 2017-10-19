package odor

// Config contains the configuration for odor.
type Config struct {
	LogLevel string `json:"logLevel" env:"LOG_LEVEL"`
}
