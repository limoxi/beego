package cron

import (
	"context"
	"fmt"
	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/orm"
	"github.com/kfchen81/beego/toolbox"
	"github.com/kfchen81/beego/vanilla"
	"math"
	"runtime/debug"
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
	ctx = context.WithValue(ctx, "orm", o)
	resource := GetManagerResource(ctx)
	ctx = context.WithValue(ctx, "jwt", resource.CustomJWTToken)
	resource.Ctx = ctx
	inst.Init(ctx, o, resource)
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

func fetchData(pi pipeInterface){
	go func(){
		defer func(){
			if err := recover(); err!=nil{
				beego.Warn(string(debug.Stack()))
				fetchData(pi)
				dingMsg := fmt.Sprintf("> goroutine from task(%s) dead \n\n 错误信息: %s \n\n", pi.(taskInterface).GetName(), err.(error).Error())
				vanilla.NewDingBot().Use("xiuer").Error(dingMsg)
			}
		}()
		for{
			data := pi.GetData()
			if data != nil{
				taskCtx := newTaskCtx()
				pi.RunConsumer(data, taskCtx)
			}
		}
	}()
}

func RegisterPipeTask(pi pipeInterface, spec string) *CronTask{
	task := RegisterTask(pi.(taskInterface), spec)
	for i := int(math.Ceil(float64(pi.GetCap())/10)); i>0; i--{
		fetchData(pi)
	}
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
		beego.Info("[cron] create cron task ", cronTask.name, cronTask.spec)
		task := toolbox.NewTask(cronTask.name, cronTask.spec, cronTask.taskFunc)
		toolbox.AddTask(cronTask.name, task)
	} else {
		for _, cronTask := range name2task {
			beego.Info("[cron] create cron task ", cronTask.name, cronTask.spec)
			task := toolbox.NewTask(cronTask.name, cronTask.spec, cronTask.taskFunc)
			toolbox.AddTask(cronTask.name, task)
		}
	}

	toolbox.StartTask()
}
