package nsq

import (
	"github.com/nsqio/go-nsq"
)

type MqReceiver interface {
	AddTask(task ...MqTask) MqReceiver
	Start()
}

type NsqReceiver struct {
	baseHost *MqHostConfigs
	Index    int8
	tasks    []MqTask
}

func NewMqReceiver(config *MqHostConfigs) MqReceiver {
	config.IsValid()
	return &NsqReceiver{baseHost: config}
}

func (r *NsqReceiver) AddTask(tasks ...MqTask) MqReceiver {
	for _, task := range tasks {
		r.tasks = append(r.tasks, task)
	}
	return r
}

func (r *NsqReceiver) Start() {
	if len(r.tasks) <= 0 {
		panic("add task first")
	}
	for r.Index < int8(len(r.tasks)) {
		r.tasks[r.Index].run(r.baseHost)
		r.Index ++
	}
	// 当前运行任务数量即为 r.Index
}

type MqTask interface {
	run(config *MqHostConfigs)
}

type MqHostConfigs struct {
	Lookup, Nsq []string
}

func (c *MqHostConfigs) IsValid() {
	if c == nil {
		panic("mqTaskConfigs is nil")
	}
	if len(c.Nsq) <= 0 && len(c.Lookup) <= 0 {
		panic("need at least one of Nsq or Lookup addrs")
	}
}

func NewNsqTask(topic, channel string, h nsq.Handler, optionalHost ...*MqHostConfigs) MqTask {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	consumer.AddHandler(h)
	task := &NsqTask{
		consumer: consumer,
		Fatal:    err,
	}
	if len(optionalHost) > 0 && optionalHost[0] != nil {
		task.OptionalHost = optionalHost[0]
	}
	return task
}

type NsqTask struct {
	consumer     *nsq.Consumer
	OptionalHost *MqHostConfigs
	ConErr       []error
	Fatal        error
}

func (task *NsqTask) run(config *MqHostConfigs) {
	if task.Fatal != nil {
		return
	}
	if task.OptionalHost != nil {
		if err := task.consumer.ConnectToNSQLookupds(task.OptionalHost.Lookup); err != nil {
			task.ConErr = append(task.ConErr, err)
		}
		if err := task.consumer.ConnectToNSQDs(task.OptionalHost.Nsq); err != nil {
			task.ConErr = append(task.ConErr, err)
		}
		return
	}
	if err := task.consumer.ConnectToNSQLookupds(config.Lookup); err != nil {
		task.ConErr = append(task.ConErr, err)
	}
	if err := task.consumer.ConnectToNSQDs(config.Nsq); err != nil {
		task.ConErr = append(task.ConErr, err)
	}
}
