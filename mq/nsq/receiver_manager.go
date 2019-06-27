package nsq

import "github.com/nsqio/go-nsq"

var log   logger
var logLv nsq.LogLevel

func init() {
	//log =
	logLv = nsq.LogLevelDebug
}

func SetLogLv(lv nsq.LogLevel) {
	logLv = lv
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
		//config.log =
		//config.logLv = logLv
		task := NewNsqTask()
		task.set(config)
		a.receiver.AddTask(task)
	}
}

func (a *ReceiverManager) Start() {
	a.receiver.Start()
}