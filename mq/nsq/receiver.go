package nsq

import (
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

type MqReceiver interface {
	AddTask(task ...MqTask) MqReceiver
	BaseHost() (baseHost *MqHostConfigs)
	Start()
	Stop()
}

type NsqReceiver struct {
	baseHost *MqHostConfigs
	Index    uint32
	tasks    []MqTask
}

func NewNsqReceiver(config *MqHostConfigs) MqReceiver {
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

func (r *NsqReceiver) retry(index uint32, duration time.Duration) {
	for !r.tasks[index].connected() {
		fmt.Println(fmt.Sprintf("task[%d] not connected, retry after %v", index, duration))
		time.Sleep(duration)
		r.tasks[index].run()
	}
	fmt.Println(fmt.Sprintf("retry for task[%d] success", index))
}

func (r *NsqReceiver) Start() {
	if len(r.tasks) <= 0 {
		panic("Please add a task first")
	}
	for r.Index < uint32(len(r.tasks)) {
		r.tasks[r.Index].run()
		// if not connected, retry every 5 min
		if !r.tasks[r.Index].connected() {
			go r.retry(r.Index, 5*time.Minute)
		}
		r.Index++
	}
}

func (r *NsqReceiver) Stop() {
	for _, task := range r.tasks {
		task.stop()
	}
}

type MqHostConfigs struct {
	Lookup, Nsq []string
}

type MqTaskConfigs struct {
	Topic, Channel string
	Handler        nsq.HandlerFunc
	Host           *MqHostConfigs
	Configs        map[string]interface{} // nsq tls configs
	Concurrency    int // goroutines numbers
}

type MqTask interface {
	set(config *MqTaskConfigs)
	run()
	connected() bool
	stop()
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
	nsqConf := nsq.NewConfig()
	for option, value := range config.Configs {
		if err := nsqConf.Set(option, value); err != nil {
			task.Fatal = err
			return
		}
	}
	consumer, err := nsq.NewConsumer(config.Topic, config.Channel, nsqConf)
	if err != nil {
		task.Fatal = err
		return
	}
	consumer.SetLogger(nsqLog, nsqLogLv)
	if config.Concurrency > 1 {
		consumer.AddConcurrentHandlers(config.Handler, config.Concurrency)
	} else {
		consumer.AddHandler(config.Handler)
	}
	task.consumer = consumer
	task.host = config.Host
}

func (task *NsqTask) run() {
	if task.Fatal != nil {
		panic(task.Fatal)
	}
	task.ConErr = nil

	//for _, url := range task.host.Lookup {
	//	if err := task.consumer.ConnectToNSQLookupd(url); err != nil {
	//		task.ConErr = append(task.ConErr, err)
	//	}
	//}

	for _, url := range task.host.Nsq {
		if err := task.consumer.ConnectToNSQD(url); err != nil {
			task.ConErr = append(task.ConErr, err)
		}
	}
}

func (task *NsqTask) connected() bool {
	totalSetups := len(task.host.Lookup) + len(task.host.Nsq)
	if len(task.ConErr) < totalSetups || totalSetups == 0 {
		return true
	}
	for _, err := range task.ConErr {
		if err == nil {
			return true
		}
	}
	return false
}

func (task *NsqTask) stop() {
	task.consumer.Stop()
}
