package cron

import (
	"context"
	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/orm"
	"github.com/kfchen81/beego/toolbox"
	"github.com/kfchen81/beego/vanilla"
)

type CronTask struct {
	name string
	spec string
	taskFunc toolbox.TaskFunc
	onlyRunThisTask bool
}

func (this *CronTask) OnlyRun() {
	this.onlyRunThisTask = true
}

var name2task = make(map[string]*CronTask)

func newTaskCtx() *TaskContext{
	inst := new(TaskContext)
	ctx := context.Background()
	o := orm.NewOrm()
	resource := GetManagerResource(ctx)
	inst.SetCtx(ctx, o).SetResource(resource)
	return inst
}

func taskWrapper(task taskInterface) toolbox.TaskFunc{

	return func() error{
		taskCtx := newTaskCtx()
		o := taskCtx.GetOrm()
		ctx := taskCtx.GetCtx()

		defer vanilla.RecoverFromCronTaskPanic(ctx)
		if task.IsEnableTx(){
			o.Begin()
			fnErr := task.Run(taskCtx)
			o.Commit()
			return fnErr
		}else{
			return task.Run(taskCtx)
		}
	}
}

func RegisterPipeTask(pi pipeInterface, spec string) *CronTask{
	task := RegisterTask(pi.(taskInterface), spec)
	pi.RunConsumer()
	return task
}

func RegisterTask(task taskInterface, spec string) *CronTask {
	tname := task.GetName()
	wrappedFn := taskWrapper(task)
	cronTask := &CronTask{
		name: tname,
		spec: spec,
		taskFunc: wrappedFn,
		onlyRunThisTask: false,
	}
	name2task[tname] = cronTask

	return cronTask
}

func RegisterCronTask(tname string, spec string, f toolbox.TaskFunc) *CronTask {
	cronTask := &CronTask{
		name: tname,
		spec: spec,
		taskFunc: f,
		onlyRunThisTask: false,
	}
	name2task[tname] = cronTask

	return cronTask
}

func StartCronTasks() {
	var onlyRunTask *CronTask
	for _, task := range name2task {
		if task.onlyRunThisTask {
			onlyRunTask = task
		}
	}

	if onlyRunTask != nil {
		cronTask := onlyRunTask
		beego.Info("[cron] create cron task ", cronTask.name)
		task := toolbox.NewTask(cronTask.name, cronTask.spec, cronTask.taskFunc)
		toolbox.AddTask(cronTask.name, task)
	} else {
		for _, cronTask := range name2task {
			beego.Info("[cron] create cron task ", cronTask.name)
			task := toolbox.NewTask(cronTask.name, cronTask.spec, cronTask.taskFunc)
			toolbox.AddTask(cronTask.name, task)
		}
	}

	toolbox.StartTask()
}
