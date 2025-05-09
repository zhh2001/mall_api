package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	//sugar := logger.Sugar()
	zap.ReplaceGlobals(logger)
}
