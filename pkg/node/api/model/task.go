// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package model

import (
	"time"

	"github.com/pquerna/ffjson/ffjson"
)

//Shell 执行脚本配置
type Shell struct {
	Cmd []string `json:"cmd"`
}

//TaskTemp 任务模版
type TaskTemp struct {
	Name    string            `json:"name" validate:"name|required"`
	ID      string            `json:"id" validate:"id|uuid"`
	Shell   Shell             `json:"shell"`
	Envs    map[string]string `json:"envs"`
	Input   string            `json:"input"`
	Args    []string          `json:"args"`
	Depends []string          `json:"depends"`
	Timeout int               `json:"timeout|required|numeric"`
	//OutPutChan
	//结果输出通道，错误输出OR标准输出
	OutPutChan string            `json:"out_put_chan" validate:"out_put_chan|required|in:stdout,stderr"`
	CreateTime time.Time         `json:"create_time"`
	Labels     map[string]string `json:"labels"`
}

func (t TaskTemp) String() string {
	res, _ := ffjson.Marshal(&t)
	return string(res)
}

//Task 任务
type Task struct {
	Name   string    `json:"name" validate:"name|required"`
	ID     string    `json:"id" validate:"id|uuid"`
	TempID string    `json:"temp_id,omitempty" validate:"temp_id|uuid"`
	Temp   *TaskTemp `json:"temp,omitempty"`
	//执行的节点
	Nodes []string `json:"nodes"`
	//执行时间定义
	//例如每30分钟执行一次:@every 30m
	Timer   string `json:"timer"`
	TimeOut int64  `json:"time_out"`
	// 执行任务失败重试次数
	// 默认为 0，不重试
	Retry int `json:"retry"`
	// 执行任务失败重试时间间隔
	// 单位秒，如果不大于 0 则马上重试
	Interval int `json:"interval"`
	//每个执行节点执行状态
	Status       map[string]TaskStatus `json:"status,omitempty"`
	CreateTime   time.Time             `json:"create_time"`
	StartTime    time.Time             `json:"start_time"`
	CompleteTime time.Time             `json:"complete_time"`
	ResultPath   string                `json:"result_path"`
	EventID      string                `json:"event_id"`
	IsOnce       bool                  `json:"is_once"`
	OutPut       []*TaskOutPut         `json:"out_put"`
}

func (t Task) String() string {
	res, _ := ffjson.Marshal(&t)
	return string(res)
}

//CanBeDelete 能否被删除
func (t Task) CanBeDelete() bool {
	if t.Status == nil || len(t.Status) == 0 {
		return true
	}
	for _, v := range t.Status {
		if v.Status == "exec" {
			return false
		}
	}
	return true
}

//TaskOutPut 任务输出
type TaskOutPut struct {
	NodeID string            `json:"node_id"`
	Global map[string]string `json:"global"`
	Inner  map[string]string `json:"inner"`
	//返回数据类型，检测结果类(check) 执行安装类 (install) 普通类 (common)
	Type   string             `json:"type"`
	Status []TaskOutPutStatus `json:"status"`
}

//ParseTaskOutPut json parse
func ParseTaskOutPut(body string) (t TaskOutPut, err error) {
	err = ffjson.Unmarshal([]byte(body), &t)
	return
}

//TaskOutPutStatus 输出数据
type TaskOutPutStatus struct {
	Name            string `json:"name"`
	ConditionType   string `json:"condition_type"`
	ConditionStatus string `json:"condition_status"`
}

//TaskStatus 任务状态
type TaskStatus struct {
	Status       string    `json:"status"` //执行状态，create init exec complete timeout
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	TakeTime     int       `json:"take_time"`
	CompleStatus string    `json:"comple_status"`
	//脚本退出码
	ShellCode int `json:"shell_code"`
}

//TaskGroup 任务组
type TaskGroup struct {
	Name       string           `json:"name" validate:"name|required"`
	ID         string           `json:"id" validate:"id|uuid"`
	Tasks      []*Task          `json:"tasks"`
	CreateTime time.Time        `json:"create_time"`
	Status     *TaskGroupStatus `json:"status"`
}

func (t TaskGroup) String() string {
	res, _ := ffjson.Marshal(&t)
	return string(res)
}

//CanBeDelete 是否能被删除
func (t TaskGroup) CanBeDelete() bool {
	if t.Status == nil || len(t.Status.TaskStatus) == 0 {
		return true
	}
	for _, v := range t.Status.TaskStatus {
		if v.Status == "exec" {
			return false
		}
	}
	return true
}

//TaskGroupStatus 任务组状态
type TaskGroupStatus struct {
	TaskStatus map[string]TaskStatus `json:"task_status"`
	InitTime   time.Time             `json:"init_time"`
	StartTime  time.Time             `json:"start_time"`
	EndTime    time.Time             `json:"end_time"`
	Status     string                `json:"status"` //create init exec complete timeout
}
