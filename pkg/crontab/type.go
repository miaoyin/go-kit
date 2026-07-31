package crontab

import "time"

//
//任务编排方法
//  方法1. module+hanlder反射
//  方法2. 顶层维护map[string]func
//

//-------------------------------
//  配置
//-------------------------------

//Configuration 配置参数
// 1.推荐的的配置
// 2.支持泛型自定义配置
type Configuration struct{
	Enable bool
	Tasks []CronTask
}

//CronTask 定时任务
type CronTask struct{
	//任务名
	Name       string
	//定时描述
	Spec       string
	//启用
	Enable     bool     `json:",omitempty"`
	//模块名
	Module     string   `json:",omitempty"`
	//任务函数, 用于反射获取函数
	Handler    string  `json:",omitempty"`
}

func (c *CronTask) GetName() string{
	return c.Name
}

func (c *CronTask) GetSpec() string{
	return c.Spec
}


//-------------------------------
//  任务信息
//-------------------------------

//ManagedTask 任务信息
//  记录Name信息
type ManagedTask struct{
	//任务信息
	Task
	//cron信息
	ID    int
	//上次运行时间
	Prev   time.Time
	//下次运行时间
	Next  time.Time
}


//Task 任务接口
type Task interface{
	GetName() string
	GetSpec() string
}
