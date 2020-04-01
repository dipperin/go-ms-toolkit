package scheduler

import (
	"errors"
	"github.com/dipperin/go-ms-toolkit/log"
	"github.com/dipperin/go-ms-toolkit/qyenv"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"strconv"
)

type task []*QyTask

// 新建任务
func NewScheduler() *task {
	t := make([]*QyTask, 0)
	ta := task(t)
	return &ta
}

// 添加任务
// ps: 参数为字符串则按照cron规则解析, 为数字则按照 每n分钟解析
func (task *task) Add(name string, prodSpecOrPerMin interface{}, f func(), nonProdSpecOrPerMin ...interface{}) *task {
	var spec string
	switch prodSpecOrPerMin.(type) {
	case string:
		spec = prodSpecOrPerMin.(string)
	case int:
		spec = "0 0/" + strconv.Itoa(absInt(prodSpecOrPerMin.(int))) + " * * * ?"
	default:
		panic("生产配置类型错误, 请传string 或 int !!!")
	}
	// 判断非生产情况
	if len(nonProdSpecOrPerMin) > 0 && qyenv.GetUseDocker() != 2 {
		switch nonProdSpecOrPerMin[0].(type) {
		case string:
			*task = append(*task, &QyTask{name, f, nonProdSpecOrPerMin[0].(string)})
		case int:
			minute := nonProdSpecOrPerMin[0].(int)
			// 如果是数字, 取绝对值
			tmp := "0 0/" + strconv.Itoa(absInt(minute)) + " * * * ?"
			*task = append(*task, &QyTask{name, f, tmp})
		default:
			log.QyLogger.Warn("QySchedule  task.Add(nonProdSpecOrPerMin), param is not valid, task set to default spec.")
			*task = append(*task, &QyTask{name, f, spec})
		}
		return task
	}
	*task = append(*task, &QyTask{name, f, spec})
	return task
}

// 启动服务
func (task *task) Run(async bool) error {
	log.QyLogger.Info("current task number:", zap.Int("number", len(*task)))
	return startScheduler(*task, async)
}

func absInt(a int) int {
	if a > 0 {
		return a
	}
	return a * -1
}


type QyTask struct {
	// 任务名称
	Name string
	// 需要执行的任务
	Run func() ()
	// 任务执行的周期配置：如,每天凌晨1点执行一次：0 0 1 * * ?
	Spec string
}

// 启动启元
func startScheduler(tasks []*QyTask, async bool) error {
	if len(tasks) == 0 {
		return errors.New("start timer task, but no task")
	}
	c := cron.New()
	for _, t := range tasks {
		log.QyLogger.Info("start timer task:", zap.String("name", t.Name), zap.String("time config", t.Spec))
		if err := c.AddFunc(t.Spec, t.Run); err != nil {
			return err
		}
	}
	if async {
		c.Start()
	} else {
		c.Run()
	}
	return nil
}