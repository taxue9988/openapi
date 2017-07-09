package service

import (
	"encoding/json"
	"log"

	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func initLogger(lp string, lv string, isDebug bool) {
	js := fmt.Sprintf(`{
		"level": "%s",
		"encoding": "json",
		"outputPaths": ["stdout","%s"],
		"errorOutputPaths": ["stderr","%s"]
	}`, lv, lp, lp)

	rawJSON := []byte(js)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	var err error
	Logger, err = cfg.Build()
	if err != nil {
		log.Fatal("init logger error: ", err)
	}

	// n := 0
	// for {
	// 	Logger.Info("logger初始化成功", zap.Int("n", n))
	// 	n++

	// 	time.Sleep(5 * time.Second)
	// }
	Logger.Info("logger初始化成功")
}
