package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

var QyLogger *zap.Logger

func init() {
	if QyLogger == nil {
		cfg := zap.NewDevelopmentConfig()
		// set log output
		cfg.OutputPaths = []string{"stdout"}
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		QyLogger, _ = cfg.Build()
	}
}

// init logger
func InitLogger(lvl zapcore.Level, targetDir, logFileName string, withConsole bool) {

	logger, err := newLogConfig(lvl, targetDir, logFileName, withConsole).Build()

	if err != nil {
		panic(err)
	}

	QyLogger = logger

}

func LoggerEnd() {
	if QyLogger == nil {
		return
	}

	_ = QyLogger.Sync()
}

// new log config
func newLogConfig(lvl zapcore.Level, targetDir, logFileName string, withConsole bool) *zap.Config {
	out := getOutPaths(targetDir, logFileName)
	if withConsole {
		out = append(out, "stdout")
	}

	eConfig := zap.NewProductionEncoderConfig()
	eConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	eConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	config := &zap.Config{
		Level:       zap.NewAtomicLevelAt(lvl),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    eConfig,
		OutputPaths:      out,
		ErrorOutputPaths: []string{"stderr"},
	}

	return config
}

func getOutPaths(targetDir, logFileName string) (out []string) {
	// Returns the default if path is empty
	if logFileName == "" {
		return
	}

	if !pathExists(targetDir) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	out = append(out, filepath.Join(targetDir, logFileName))

	return
}

// Determine if the path file exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
