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
		cfg.DisableCaller = true
		// set log output
		cfg.OutputPaths = []string{"stdout"}
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		QyLogger, _ = cfg.Build()
	}
}

// init logger
func InitLogger(lvl zapcore.Level, targetDir, logFileName string, withConsole bool) {
	QyLogger = zap.New(
		newLogCore(lvl, targetDir, logFileName, withConsole),
		newLogOptions()...,
	)
}

func LoggerEnd() {
	if QyLogger == nil {
		return
	}

	_ = QyLogger.Sync()
}

// new log core
func newLogCore(lvl zapcore.Level, targetDir, logFileName string, withConsole bool) zapcore.Core {
	out := getOutPaths(targetDir, logFileName)
	if withConsole {
		out = append(out, "stdout")
	}
	if logFileName != "" {
		logFileName = "err_" + logFileName
	}
	errOut := getOutPaths(targetDir, logFileName)
	errOut = append(errOut, "stderr")
	sink, _, err := zap.Open(out...)
	if err != nil {
		panic("outputPaths open err, err="+err.Error())
	}
	errSink, _, err := zap.Open(errOut...)
	if err != nil {
		panic("errOutputPaths open err, err="+err.Error())
	}

	eConfig := zap.NewProductionEncoderConfig()
	eConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	eConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(eConfig)


	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, sink, normalLevelEnable{flagLevel: lvl}),
		zapcore.NewCore(consoleEncoder, errSink,
			zap.LevelEnablerFunc(func(zl zapcore.Level) bool {
			return zl >= zapcore.ErrorLevel
		})),
	)
}

type normalLevelEnable struct {
	flagLevel zapcore.Level
}

func (c normalLevelEnable) Enabled(lvl zapcore.Level) bool {
	return lvl >= c.flagLevel  && lvl < zap.ErrorLevel
}

func newLogOptions() []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
}

func getOutPaths(targetDir, logFileName string) (out []string) {
	// Returns the default if path is empty
	if logFileName == "" {
		return
	}

	if !pathExists(targetDir) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			panic(err.Error()+"; dir="+targetDir)
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
