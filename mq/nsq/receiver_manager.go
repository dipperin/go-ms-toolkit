package nsq

import (
	"github.com/nsqio/go-nsq"
	"log"
	"os"
)

var nsqLog   logger
var nsqLogLv nsq.LogLevel

func init() {
	nsqLog =  log.New(os.Stderr, "", log.LstdFlags)
	nsqLogLv = nsq.LogLevelInfo
}

func SetLog(l logger) {
	nsqLog = l
}

func SetLogLv(lv nsq.LogLevel) {
	nsqLogLv = lv
}

type NsqHandler interface {
	TaskConfig() (config *MqTaskConfigs)
}

type NsqHandlerFunc func() (config *MqTaskConfigs)

func (f NsqHandlerFunc) TaskConfig() (config *MqTaskConfigs) {
	return f()
}

type ReceiverManager struct {
	receiver MqReceiver
}

func NewReceiverManager(receiver MqReceiver, h ...NsqHandler) *ReceiverManager {
	rm := &ReceiverManager{ receiver: receiver }
	rm.Add(h...)
	return rm
}

func (a *ReceiverManager) Add(h ...NsqHandler) {
	for i := range h {
		config := h[i].TaskConfig()
		if config.Host == nil { // default base Host
			config.Host = a.receiver.BaseHost()
		}
		task := NewNsqTask()
		task.set(config)
		a.receiver.AddTask(task)
	}
}

func (a *ReceiverManager) Start() {
	a.receiver.Start()
}

func (a *ReceiverManager) Stop() {
	a.receiver.Stop()
}