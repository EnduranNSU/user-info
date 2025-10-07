package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/num30/config"
	"github.com/rs/zerolog/log"

	appconfig "github.com/EnduranNSU/end-user-info/config"
	"github.com/EnduranNSU/end-user-info/internal/logging"
)

func init() {
	// Setup default logger
	logging.SetupLogger(
		logging.Config{
			Level: "info",
			Console: logging.ConsoleLoggerConfig{
				Enable:   true,
				Encoding: "text",
			},
			File: logging.FileLoggerConfig{
				Enable: false,
			},
		},
	)
}

func main() {
	// Load config
	var cfg appconfig.Config
	configName := appconfig.GetConfigName()

	err := config.NewConfReader(configName).WithPrefix("APP").Read(&cfg)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to load config")
	}

	// Setup logger
	logging.SetupLogger(toLoggerConfig(cfg.Logger))
}

func toLoggerConfig(cfg appconfig.LoggerConfig) logging.Config {
	return logging.Config{
		Level: cfg.Level,
		Console: logging.ConsoleLoggerConfig{
			Enable:   cfg.Console.Enable,
			Encoding: cfg.Console.Encoding,
		},
		File: logging.FileLoggerConfig{
			Enable:  cfg.File.Enable,
			DirPath: cfg.File.DirPath,
			MaxSize: cfg.File.MaxSize,
			MaxAge:  cfg.File.MaxAge,
		},
	}
}
