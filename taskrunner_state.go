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

func (tr *taskRunnerState) setState(runner TaskRunner, state TaskState) {
	atomic.StoreInt32(&tr.state, int32(state))
	if tr.notifier != nil {
		tr.notifier <- runner
	}
}
