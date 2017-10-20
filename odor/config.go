package odor

// Config contains the configuration for odor.
type Config struct {
	LogLevel string              `json:"logLevel" env:"LOG_LEVEL"`
	Address  string              `json:"address" env:"ADDRESS"`
	Filters  map[string][]string `json:"filters"`
}
