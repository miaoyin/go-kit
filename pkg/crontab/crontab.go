package crontab

import (
	"context"
	"fmt"
	"sync"

	"github.com/miaoyin/go-kit/pkg/module"
	"github.com/robfig/cron/v3"
	"golang.org/x/exp/maps"
)

//支持Module接口
var _ module.Module[module.EmptyConfig] = &Crontab[module.EmptyConfig]{}


func NewCrontab[T any](name string, initCfg T)*Crontab[T]{
	c := Crontab[T]{
		cron: cron.New(),
		BaseModule: module.NewBaseModule[T](name, initCfg),
		nameIndex: make(map[string]*ManagedTask),
		idIndex: make(map[int]*ManagedTask),
	}
	_ = c.Start(nil)
	return &c
}

//BuildSimpleCrontab
//  1. name默认crontab
//  2. 忽略配置信息
func BuildSimpleCrontab() *Crontab[module.EmptyConfig]{
	return NewCrontab[module.EmptyConfig]("crontab", module.EmptyConfig{})
}

//Crontab 定时任务
// 用途
//  1. 实现Module接口
//  2. 缓存任务信息, cron默认只有id
type Crontab[T any] struct{
	//模块接口
	*module.BaseModule[T]

	//任务调度
	cron 	*cron.Cron

	//缓存信息
	nameIndex map[string]*ManagedTask
	idIndex map[int]*ManagedTask
	mu  sync.RWMutex
}

func (ct *Crontab[T]) Start(_ context.Context) error{
	ct.mu.Lock()
	defer ct.mu.Unlock()
	if ct.IsRunning(){
		return nil
	}
	ct.cron.Start()
	ct.DoStart()
	return nil
}

func (ct *Crontab[T]) Close(_ context.Context) error{
	ct.mu.Lock()
	defer ct.mu.Unlock()
	if !ct.IsRunning(){
		return nil
	}
	ct.cron.Stop()
	//清空
	ct.nameIndex = make(map[string]*ManagedTask)
	ct.idIndex = make(map[int]*ManagedTask)
	ct.DoClose()
	return nil
}

func (ct *Crontab[T]) RegisterSimpleTask(name string, spec string, cmd func()) (int, error){
	return ct.RegisterTask(&CronTask{Name: name, Spec: spec}, cmd)
}

//RegisterTask 增加任务
func (ct *Crontab[T]) RegisterTask(task Task, cmd func()) (int, error){
	ct.mu.Lock()
	defer ct.mu.Unlock()
	taskName := task.GetName()
	// 重名校验
	if _, exists := ct.nameIndex[taskName]; exists {
		return 0, fmt.Errorf("task name %s already exists", taskName)
	}

	//添加任务
	id, err := ct.cron.AddFunc(task.GetSpec(), cmd)
	if err != nil {return 0, err}
	//存cache
	intId := int(id)
	mgrTask := &ManagedTask{
		ID: intId,
		Task: task,
	}
	ct.nameIndex[taskName] = mgrTask
	ct.idIndex[intId] = mgrTask
	return intId, nil
}


//Remove 删除任务
func (ct *Crontab[T]) Remove(id int){
	ct.mu.Lock()
	defer ct.mu.Unlock()
	task, ok := ct.idIndex[id]
	if ok{
		delete(ct.idIndex, id)
		delete(ct.nameIndex, task.GetName())
	}
	ct.cron.Remove(cron.EntryID(id))
}

//List 任务列表
func (ct *Crontab[T]) List() []*ManagedTask{
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	jobs := make([]*ManagedTask, 0)
	for _, item := range ct.cron.Entries(){
		task := ct.idIndex[int(item.ID)]
		jobs = append(jobs, &ManagedTask{
			Task:task,
			ID:int(item.ID),
			Prev:item.Prev,
			Next:item.Next,
		})
	}
	return jobs
}

func (ct *Crontab[T]) Map() map[string]*ManagedTask{
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	copyTask := make(map[string]*ManagedTask)
	maps.Copy(copyTask, ct.nameIndex)
	return copyTask
}

func (ct *Crontab[T]) Cron() *cron.Cron{
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	return ct.cron
}
