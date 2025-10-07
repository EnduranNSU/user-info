package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/num30/config"
	"github.com/rs/zerolog/log"

	appconfig "github.com/EnduranNSU/end-user-info/config"
	"github.com/EnduranNSU/end-user-info/internal/db/repository/impl"
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

	//Open db
	db, err := sql.Open("postgres", 
	fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable host=%s port=%d",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Dbname, cfg.Db.Host, cfg.Db.Port))
    if err != nil {
        log.Fatal().Stack().Err(err).Msgf("Failed to connect to database: %v", err)
    }
    defer db.Close()

	// Checking connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil {
        log.Fatal().Stack().Err(err).Msgf("Failed to ping database: %v", err)
    }

	//init repo
	repo := impl.NewUserRepository(db)
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
