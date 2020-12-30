// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package executor

import "sync/atomic"

type taskRunnerState struct {
	notifier chan<- TaskState
	state    int32
}

func (tr *taskRunnerState) SetNotifier(notifier chan<- TaskState) {
	tr.notifier = notifier
}

//是否有任务正在执行
func (tr *taskRunnerState) IsRunning() bool {
	return atomic.LoadInt32(&tr.state) == TaskStateRunning
}

func (tr *taskRunnerState) setState(state TaskState) {
	atomic.StoreInt32(&tr.state, int32(state))
	if tr.notifier != nil {
		tr.notifier <- state
	}
}
