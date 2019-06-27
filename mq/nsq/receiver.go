package nsq

import (
	"github.com/nsqio/go-nsq"
)

type MqReceiver interface {
	AddTask(task ...MqTask) MqReceiver
	BaseHost() (baseHost *MqHostConfigs)
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

func (r *NsqReceiver) BaseHost() (baseHost *MqHostConfigs) {
	return r.baseHost
}

func (r *NsqReceiver) Start() {
	if len(r.tasks) <= 0 {
		panic("Please add a task first")
	}
	for r.Index < int8(len(r.tasks)) {
		r.tasks[r.Index].run()
		r.Index ++
	}
	// current task number = r.Index
}

type MqHostConfigs struct {
	Lookup, Nsq []string
}

type MqTaskConfigs struct {
	topic, channel string
	handler        nsq.HandlerFunc
	host           *MqHostConfigs
}

type MqTask interface {
	set(config *MqTaskConfigs)
	run()
}

func (c *MqHostConfigs) IsValid() {
	if c == nil {
		panic("mqTaskConfigs is nil")
	}
	if len(c.Nsq) <= 0 && len(c.Lookup) <= 0 {
		panic("need at least one of Nsq or Lookup addrs")
	}
}

func NewNsqTask() MqTask {
	return &NsqTask{}
}

type NsqTask struct {
	consumer *nsq.Consumer
	host     *MqHostConfigs
	Fatal    error
	ConErr   []error
}

func (task *NsqTask) set(config *MqTaskConfigs) {
	consumer, err := nsq.NewConsumer(config.topic, config.channel, nsq.NewConfig())
	if err != nil {
		task.Fatal = err
		return
	}
	consumer.AddHandler(config.handler)
	task.consumer = consumer
	task.host = config.host
}

func (task *NsqTask) run() {
	if task.Fatal != nil {
		panic(task.Fatal)
	}

	for _, url := range task.host.Lookup {
		if err := task.consumer.ConnectToNSQLookupd(url); err != nil {
			task.ConErr = append(task.ConErr, err)
		}
	}

	for _, url := range task.host.Nsq {
		if err := task.consumer.ConnectToNSQD(url); err != nil {
			task.ConErr = append(task.ConErr, err)
		}
	}
}
