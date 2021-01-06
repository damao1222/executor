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

type TaskRunnerFIFO struct {
	task chan Task
	stop chan struct{}
	once sync.Once

	taskRunnerState
}

func NewFIFO(taskSize int) *TaskRunnerFIFO {
	return &TaskRunnerFIFO{
		task:            make(chan Task, taskSize),
		stop:            make(chan struct{}),
		taskRunnerState: taskRunnerState{nil, TaskStateIdle},
	}
}

func (tr *TaskRunnerFIFO) SetTask(task Task) bool {
	select {
	case tr.task <- task:
		return true
	default:
		return false
	}
}

func (tr *TaskRunnerFIFO) Stop() {
	tr.once.Do(func() {
		close(tr.stop)
	})
}

func (tr *TaskRunnerFIFO) Next() {

}

func (tr *TaskRunnerFIFO) OnExpired(task Task) {
}

func (tr *TaskRunnerFIFO) Steal(runner TaskRunner) bool {
	// not idle
	if !tr.isIdle() {
		return false
	}
	if other, ok := runner.(*TaskRunnerFIFO); ok {
		select {
		case task, ok := <-other.task:
			if ok {
				tr.SetTask(task)
			}
			return true
		default:
			return false
		}
	} else {
		return false
	}
}

//是否有任务正在执行
func (tr *TaskRunnerFIFO) IsIdle() bool {
	return tr.isIdle() && len(tr.task) == 0
}

func (tr *TaskRunnerFIFO) Loop() {
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

func (tr *TaskRunnerFIFO) handlePanic() {
	if r := recover(); r != nil {
		log.Print("task panic: ", r)
	}
}

func (tr *TaskRunnerFIFO) safeRun(task Task) {
	defer tr.handlePanic()

	tr.setRunning(tr)
	defer tr.setIdle(tr)

	task()
}
