package config

type Config struct {
	DSN string
	Addr       string
	Production bool
	Key        string
}

func NewConfig() *Config {
	return &Config{
		Addr:       ":8080",
		Production: false,
		Key:        "example",
	}
}
