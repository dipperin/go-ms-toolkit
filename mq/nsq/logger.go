package nsq

import (
	"github.com/nsqio/go-nsq"
	"log"
	"os"
)

var nsqLog logger
var nsqLogLv nsq.LogLevel

func init() {
	nsqLog = log.New(os.Stdout, "[nsq-log]", log.LstdFlags)
	nsqLogLv = nsq.LogLevelInfo
}

func SetLog(l logger) {
	nsqLog = l
}

func SetLogLv(lv nsq.LogLevel) {
	nsqLogLv = lv
}