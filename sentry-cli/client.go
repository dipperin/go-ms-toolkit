package sentry_cli

import (
	"github.com/dipperin/go-ms-toolkit/env"
	"github.com/dipperin/go-ms-toolkit/log"
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
	c, err := raven.New(dsn)
	if err != nil {
		if env.GetUseDocker() == 2 {
			panic(err)
		} else {
			log.QyLogger.Info("init sentry client failed", zap.Error(err))
		}
	}
	client = c
	return client
}
