/**
 * Copyright (C) 2018-2020, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/16
 * @time 16:49
 * @version V1.0
 * Description:
 */

package executor

import (
	"log"
	"sync"
)

type TaskRunnerOnce struct {
	task chan Task
	stop chan struct{}
	once sync.Once

	taskRunnerState
}

func NewOnce() *TaskRunnerOnce {
	return &TaskRunnerOnce{
		task:            make(chan Task),
		stop:            make(chan struct{}),
		taskRunnerState: taskRunnerState{nil, TaskStateIdle},
	}
}

//NOTICE:当Loop协程没有就绪，则会一直返回false
func (tr *TaskRunnerOnce) SetTask(task Task) bool {
	select {
	case tr.task <- task:
		return true
	default:
		return false
	}
}

func (tr *TaskRunnerOnce) Stop() {
	tr.once.Do(func() {
		close(tr.stop)
	})
}

func (tr *TaskRunnerOnce) Next() {

}

func (tr *TaskRunnerOnce) OnExpired(Task) {

}

//是否有任务正在执行
func (tr *TaskRunnerOnce) IsIdle() bool {
	return tr.isIdle()
}

func (tr *TaskRunnerOnce) Loop() {
	for {
		select {
		case task, ok := <-tr.task:
			if ok {
				tr.safeRun(task)
			}
		case <-tr.stop:
			return
		}
	}
}

// cannot steal
func (tr *TaskRunnerOnce) Steal(task TaskRunner) bool {
	return false
}

func (tr *TaskRunnerOnce) handlePanic() {
	if r := recover(); r != nil {
		log.Print("task panic: ", r)
	}
}

func (tr *TaskRunnerOnce) safeRun(task Task) {
	defer tr.handlePanic()

	tr.setRunning(tr)
	defer tr.setIdle(tr)

	task()
}
