package sentry_cli

import (
	"github.com/dipperin/go-ms-toolkit/log"
	"github.com/dipperin/go-ms-toolkit/qyenv"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
	"os"
)

var client *raven.Client = nil

func Client() *raven.Client {
	if client != nil {
		return client
	}

	dsn := os.Getenv("sentry_dsn")
	log.QyLogger.Info("sentry dsn: " + dsn)
	c, err := raven.New(dsn)
	if err != nil {
		if qyenv.GetUseDocker() == 2 {
			log.QyLogger.Warn("init sentry client failed", zap.Error(err))
			return c
		} else {
			log.QyLogger.Info("init sentry client failed", zap.Error(err))
			return c
		}
	}
	client = c
	return client
}
