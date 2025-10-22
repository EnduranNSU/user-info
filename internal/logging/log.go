package logging

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

func MarshalLogData(data any) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}

func Error(err error, operation string, jsonData []byte, message string) {
	log.Error().
		Err(err).
		Str("operation", operation).
		RawJSON("data", jsonData).
		Msg(message)
}

func Debug(operation string, jsonData []byte, message string) {
	log.Debug().
		Str("operation", operation).
		RawJSON("data", jsonData).
		Msg(message)
}


func Info(operation string, jsonData []byte, message string) {
	log.Info().
		Str("operation", operation).
		RawJSON("data", jsonData).
		Msg(message)
}

func Warn(operation string, jsonData []byte, message string) {
	log.Warn().
		Str("operation", operation).
		RawJSON("data", jsonData).
		Msg(message)
}