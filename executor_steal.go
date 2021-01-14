// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package executor

import (
	"errors"
	"sync"
	"time"
)

const (
	DefaultStealInterval       = 1 * time.Second
	DefaultInitialGorutineSize = 8
	DefaultMaxGorutineSize     = 64
	DefaultMaxTaskBufSize      = 128
)

type WorkStealingExecutorOpt func(exec *WorkStealingExecutor)

type WorkStealingExecutor struct {
	initSize      int
	runnerSize    int
	bufSize       int
	stealInterval time.Duration
	runners       []TaskRunner
	curIndex      int
	runnerLock    sync.Mutex
	stopCh        chan struct{}
}

//创建一个固定大小的协程池
//size： 协程池大小
//taskBufSize： 任务缓冲大小
//需要注意的是，可以成功加入协程池的任务数量为 size * taskBufSize
//并且当第一个协程池中的协程taskBufSize被填满之前，不会尝试使用协程池中的其他协程
func NewWorkStealingExecutor(opts ...WorkStealingExecutorOpt) *WorkStealingExecutor {
	ex := &WorkStealingExecutor{
		initSize:      DefaultInitialGorutineSize,
		runnerSize:    DefaultMaxGorutineSize,
		bufSize:       DefaultMaxTaskBufSize,
		stealInterval: DefaultStealInterval,
		curIndex:      0,
		stopCh:        make(chan struct{}),
	}
	for _, opt := range opts {
		opt(ex)
	}
	for i := 0; i < ex.initSize; i++ {
		ex.runners[i] = NewFIFO(ex.bufSize)
		//start runner loop
		go ex.runners[i].Loop()
	}

	go func() {
		ticker := time.NewTicker(ex.stealInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ex.stopCh:
				return
			case <-ticker.C:
				ex.doSteal()
			}
		}
	}()
	return ex
}

func (ex *WorkStealingExecutor) Run(task Task) error {
	ex.runnerLock.Lock()
	defer ex.runnerLock.Unlock()

	//runner := ex.runners[i]
	runner := ex.selectRunner()
	if runner != nil && runner.SetTask(task) {
		return nil
	}
	return errors.New("All Runners are busy. ")
}

func (ex *WorkStealingExecutor) selectRunner() TaskRunner {
	if len(ex.runners) == ex.runnerSize {
		ex.curIndex++
		ex.curIndex = ex.curIndex % len(ex.runners)

		return ex.runners[ex.curIndex]
	}
	return nil
}

func (ex *WorkStealingExecutor) doSteal() {
	ex.runnerLock.Lock()
	defer ex.runnerLock.Unlock()

	for _, runner := range ex.runners {
		if runner.IsIdle() {
			for _, other := range ex.runners {
				if runner != other && !other.IsIdle() {
					if runner.Steal(other) {
						return
					}
				}
			}
		}
	}
}

func (ex *WorkStealingExecutor) Stop() {
	close(ex.stopCh)

	ex.runnerLock.Lock()
	defer ex.runnerLock.Unlock()

	for _, runner := range ex.runners {
		runner.Stop()
	}
}

func SetInitialGorutineSize(size int) WorkStealingExecutorOpt {
	return func(exec *WorkStealingExecutor) {
		exec.initSize = size
	}
}

func SetMaxGorutineSize(size int) WorkStealingExecutorOpt {
	return func(exec *WorkStealingExecutor) {
		exec.runnerSize = size
	}
}

func SetTaskBufSize(size int) WorkStealingExecutorOpt {
	return func(exec *WorkStealingExecutor) {
		exec.bufSize = size
	}
}

func SetStealInterval(interval time.Duration) WorkStealingExecutorOpt {
	return func(exec *WorkStealingExecutor) {
		exec.stealInterval = interval
	}
}
