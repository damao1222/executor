// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package executor

import "sync/atomic"

type taskRunnerState struct {
	notifier chan<- TaskRunner
	state    int32
}

func (tr *taskRunnerState) SetNotifier(notifier chan<- TaskRunner) {
	tr.notifier = notifier
}

func (tr *taskRunnerState) setIdle(runner TaskRunner) (ret bool) {
	ret = atomic.CompareAndSwapInt32(&tr.state, TaskStateRunning, TaskStateIdle)
	if tr.notifier != nil {
		tr.notifier <- runner
	}
	return
}

func (tr *taskRunnerState) setRunning(runner TaskRunner) (ret bool) {
	ret = atomic.CompareAndSwapInt32(&tr.state, TaskStateIdle, TaskStateRunning)
	if tr.notifier != nil {
		tr.notifier <- runner
	}
	return
}

func (tr *taskRunnerState) isIdle() bool {
	return atomic.LoadInt32(&tr.state) == TaskStateIdle
}
