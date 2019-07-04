package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
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
	errOut := getErrOutPaths(targetDir, logFileName)
	if withConsole {
		out = append(out, "stdout")
		errOut = append(errOut, "stderr")
	}

	eConfig := zap.NewProductionEncoderConfig()
	eConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	eConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(eConfig)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, getSink(out), normalLevelEnable{flagLevel: lvl}),
		zapcore.NewCore(consoleEncoder, getSink(errOut),
			zap.LevelEnablerFunc(func(zl zapcore.Level) bool {
				return zl >= zapcore.ErrorLevel
			})),
	)
}

type normalLevelEnable struct {
	flagLevel zapcore.Level
}

func (c normalLevelEnable) Enabled(lvl zapcore.Level) bool {
	return lvl >= c.flagLevel && lvl < zap.ErrorLevel
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
			panic(err.Error() + "; dir=" + targetDir)
		}
	}

	out = append(out, filepath.Join(targetDir, logFileName))

	return
}

func getErrOutPaths(targetDir, logFileName string) (out []string) {
	if logFileName == "" {
		return
	}
	return getOutPaths(targetDir, "err_"+logFileName)
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

func newRollingLogWriter(filename string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500, // mb
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
	})
}

func newZapLogWriter(outputPaths []string) zapcore.WriteSyncer {
	w, _, err := zap.Open(outputPaths...)
	if err != nil {
		panic(fmt.Sprintf("outputPaths open err, err=%v, outputPaths=%v", err.Error(), outputPaths))
	}
	return w
}

func getSink(outputPaths []string) zapcore.WriteSyncer {
	if len(outputPaths) > 0 &&
		outputPaths[0] != "stdout" &&
		outputPaths[0] != "stderr" {
		return newRollingLogWriter(outputPaths[0])
	}
	return newZapLogWriter(outputPaths)
}
