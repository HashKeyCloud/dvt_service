package common

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger Init Logger
func InitLogger(filename string) {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	if _, ok := os.LookupEnv("test"); ok {
		log.Logger = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
		return
	}

	rollingFile := &lumberjack.Logger{
		Filename:   filename,
		MaxBackups: 32, // files
		MaxSize:    16, // megabytes
		MaxAge:     1,  // days
		Compress:   true,
	}

	log.Logger = zerolog.New(rollingFile).Level(zerolog.InfoLevel).With().Caller().Timestamp().Logger()
}
