package main

import (
	"go.uber.org/zap"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./test.log",
		"stderr",
		"stdout",
	}
	return cfg.Build()
}

func main() {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
	url := "https://zhh2001.github.io/"
	sugar := logger.Sugar()
	sugar.Infow(
		"failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	logger.Info(
		"failed to fetch URL",
		zap.String("url", url),
		zap.Int("nums", 3),
		zap.Duration("backoff", time.Second),
	)
}
