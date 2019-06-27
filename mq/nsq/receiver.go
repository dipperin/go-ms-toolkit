package nsq


import (
	"github.com/nsqio/go-nsq"
)

type MqReceiver interface {
	AddTask(task ...MqTask) MqReceiver
	Start()
}

type NsqReceiver struct {
	config       *MqTaskConfigs
	Index        int8
	tasks        []MqTask
}

func NewMqReceiver(config *MqTaskConfigs) MqReceiver {
	config.IsValid()
	return &NsqReceiver{config: config}
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
		r.tasks[r.Index].run(r.config)
		r.Index ++
	}
	// 当前运行任务数量即为 r.Index
}

type MqTask interface {
	run(config *MqTaskConfigs)
}

type MqTaskConfigs struct {
	Lookup, Nsq []string
}

func (c *MqTaskConfigs) IsValid() {
	if c == nil {
		panic("mqTaskConfigs is nil")
	}
	if len(c.Nsq) <= 0 && len(c.Lookup) <= 0 {
		panic("need at least one of Nsq or Lookup addrs")
	}
}

func NewNsqTask(topic, channel string, h nsq.Handler) MqTask {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	consumer.AddHandler(h)
	return &NsqTask{
		consumer: consumer,
		fatal:    err,
	}
}

type NsqTask struct {
	consumer *nsq.Consumer
	conErr   []error
	fatal    error
}

func (task *NsqTask) run(config *MqTaskConfigs) {
	if task.fatal != nil {
		return
	}
	if err := task.consumer.ConnectToNSQLookupds(config.Lookup); err != nil {
		task.conErr = append(task.conErr, err)
	}
	if err := task.consumer.ConnectToNSQDs(config.Nsq); err != nil {
		task.conErr = append(task.conErr, err)
	}
}

