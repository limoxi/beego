package cron

import (
	"github.com/kfchen81/beego/toolbox"
	"github.com/kfchen81/beego"
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
