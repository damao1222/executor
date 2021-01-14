// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package executor

import "time"

const (
	DefaultStealInterval       = 1 * time.Second
	DefaultInitialGorutineSize = 8
	DefaultMaxGorutineSize     = 64
	DefaultMaxTaskBufSize      = 128
)

type optExecutor interface {
	setInitialGorutineSize(size int)
	setMaxGorutineSize(size int)
	setTaskBufSize(size int)
	setStealInterval(duration time.Duration)
}

type ExecutorOpt func(exec optExecutor)

func SetInitialGorutineSize(size int) ExecutorOpt {
	return func(exec optExecutor) {
		exec.setInitialGorutineSize(size)
	}
}

func SetMaxGorutineSize(size int) ExecutorOpt {
	return func(exec optExecutor) {
		exec.setMaxGorutineSize(size)
	}
}

func SetTaskBufSize(size int) ExecutorOpt {
	return func(exec optExecutor) {
		exec.setTaskBufSize(size)
	}
}

func SetStealInterval(duration time.Duration) ExecutorOpt {
	return func(exec optExecutor) {
		exec.setStealInterval(duration)
	}
}
