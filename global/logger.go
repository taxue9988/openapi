package global

import (
	"encoding/json"
	"log"

	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(lp string, lv string, isDebug bool) {
	js := fmt.Sprintf(`{
		"level": "%s",
		"encoding": "json",
		"outputPaths": ["stdout","%s"],
		"errorOutputPaths": ["stderr","%s"]
	}`, lv, lp, lp)

	var cfg zap.Config
	if err := json.Unmarshal([]byte(js), &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.InitialFields = map[string]interface{}{
		"service": Conf.Common.Service,
	}

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		log.Fatal("init logger error: ", err)
	}

	Logger.Info("logger初始化成功")
}

func DebugLog(requestID string, debugOn bool, msg string, fields ...zapcore.Field) {
	if debugOn {
		Logger.Info(msg, append(fields, zap.String("rid", requestID))...)
	}
}
