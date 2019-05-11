package cron

import (
	"context"
	"github.com/kfchen81/beego"
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

func taskWrapper(ctx context.Context, fn toolbox.TaskFunc) toolbox.TaskFunc{
	o := vanilla.GetOrmFromContext(ctx)
	if o == nil{
		beego.Warn("task run without db transaction")
		return fn
	}
	return func() error{
		o.Begin()
		defer vanilla.RecoverFromCronTaskPanic(ctx)
		fnErr := fn()
		o.Commit()
		return fnErr
	}
}

func RegisterCronTaskWithContext(tname string, spec string, f toolbox.TaskFunc, ctx context.Context) *CronTask {
	ctx = context.WithValue(ctx, "taskName", tname)
	wrappedFn := taskWrapper(ctx, f)
	if wrappedFn == nil{
		beego.Error("register task [%s] failed", tname)
		return nil
	}
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
