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
	out := getOutPath(targetDir, logFileName)
	errOut := getErrOutPath(targetDir, logFileName)

	eConfig := zap.NewProductionEncoderConfig()
	eConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	eConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(eConfig)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, getOutSink(out, withConsole), normalLevelEnable{flagLevel: lvl}),
		zapcore.NewCore(consoleEncoder, getErrOutSink(errOut, withConsole),
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

func getOutPath(targetDir, logFileName string) (out string) {
	if logFileName == "" {
		return
	}

	if !pathExists(targetDir) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			panic(err.Error() + "; dir=" + targetDir)
		}
	}

	return filepath.Join(targetDir, logFileName)
}

func getErrOutPath(targetDir, logFileName string) (out string) {
	if logFileName == "" {
		return
	}
	return getOutPath(targetDir, "err_"+logFileName)
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

func getSink(outputPath string, stds []string, withConsole bool) zapcore.WriteSyncer {
	w, _, err := zap.Open(stds...)
	if err != nil {
		panic(fmt.Sprintf("default: std open err, err=%v", err.Error()))
	}
	if outputPath == "" {
		return w
	}
	rollW := zapcore.AddSync(&lumberjack.Logger{
		Filename:   outputPath,
		MaxSize:    500, // mb
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
	})
	if withConsole {
		return zap.CombineWriteSyncers([]zapcore.WriteSyncer{rollW, w}...)
	}
	return rollW
}

func getOutSink(outputPath string, withConsole bool) zapcore.WriteSyncer {
	return getSink(outputPath, []string{"stdout"}, withConsole)
}

func getErrOutSink(outputPath string, withConsole bool) zapcore.WriteSyncer {
	return getSink(outputPath, []string{"stdout", "stderr"}, withConsole)
}
