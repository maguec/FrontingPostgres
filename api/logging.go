package api

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func SetupLogging(filepath string) *zap.SugaredLogger {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"app": "fronting"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "timeKey": "ts",
      "timeEncoder": "millis",
	    "levelEncoder": "lowercase"

	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	if filepath != "" {
		cfg.OutputPaths = []string{filepath}
		cfg.ErrorOutputPaths = []string{fmt.Sprintf("%s.err", filepath), filepath}
	}

	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	sugarLogger := logger.Sugar()
	return sugarLogger
}
