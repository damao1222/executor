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

type LazyFixBUfExecutor struct {
	initSize   int
	runnerSize int
	bufSize    int
	runners    []TaskRunner
	curIndex   int
	runnerLock sync.Mutex
}

//创建一个固定大小的协程池
//size： 协程池大小
//taskBufSize： 任务缓冲大小
//需要注意的是，可以成功加入协程池的任务数量为 size * taskBufSize
//并且当第一个协程池中的协程taskBufSize被填满之前，不会尝试使用协程池中的其他协程
func NewLazyFixBUfExecutor(opts ...ExecutorOpt) *LazyFixBUfExecutor {
	ex := &LazyFixBUfExecutor{
		initSize:   DefaultInitialGorutineSize,
		runnerSize: DefaultMaxGorutineSize,
		bufSize:    DefaultMaxTaskBufSize,
		curIndex:   0,
	}
	for _, opt := range opts {
		opt(ex)
	}
	ex.runners = make([]TaskRunner, ex.initSize)
	for i := 0; i < ex.initSize; i++ {
		ex.runners[i] = NewFIFO(ex.bufSize)
		//start runner loop
		go ex.runners[i].Loop()
	}

	return ex
}

func (ex *LazyFixBUfExecutor) Run(task Task) error {
	ex.runnerLock.Lock()
	defer ex.runnerLock.Unlock()

	if ex.selectRunner(task) {
		return nil
	}
	return errors.New("Run task failed. ")
}

func (ex *LazyFixBUfExecutor) selectRunner(task Task) bool {
	for _, runner := range ex.runners {
		if runner.IsIdle() {
			return runner.SetTask(task)
		}
	}
	if len(ex.runners) == ex.runnerSize {
		ex.curIndex++
		ex.curIndex = ex.curIndex % len(ex.runners)

		return ex.runners[ex.curIndex].SetTask(task)
	} else {
		newRunner := NewFIFO(ex.bufSize)
		//start runner loop
		go newRunner.Loop()
		ex.runners = append(ex.runners, newRunner)
		return newRunner.SetTask(task)
	}
}

func (ex *LazyFixBUfExecutor) Stop() {
	ex.runnerLock.Lock()
	defer ex.runnerLock.Unlock()

	for _, runner := range ex.runners {
		runner.Stop()
	}
}

func (ex *LazyFixBUfExecutor) setInitialGorutineSize(size int) {
	ex.initSize = size
}

func (ex *LazyFixBUfExecutor) setMaxGorutineSize(size int) {
	ex.runnerSize = size
}

func (ex *LazyFixBUfExecutor) setTaskBufSize(size int) {
	ex.bufSize = size
}

func (ex *LazyFixBUfExecutor) setStealInterval(interval time.Duration) {
}

type WorkStealingExecutor struct {
	LazyFixBUfExecutor
	stealInterval time.Duration
	stopCh        chan struct{}
}

//创建一个固定大小的协程池
//size： 协程池大小
//taskBufSize： 任务缓冲大小
//需要注意的是，可以成功加入协程池的任务数量为 size * taskBufSize
//并且当第一个协程池中的协程taskBufSize被填满之前，不会尝试使用协程池中的其他协程
func NewWorkStealingExecutor(opts ...ExecutorOpt) *WorkStealingExecutor {
	ex := &WorkStealingExecutor{
		LazyFixBUfExecutor: *NewLazyFixBUfExecutor(opts...),
		stealInterval:      DefaultStealInterval,
		stopCh:             make(chan struct{}),
	}
	for _, opt := range opts {
		opt(ex)
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

	ex.LazyFixBUfExecutor.Stop()
}

func (ex *WorkStealingExecutor) setStealInterval(interval time.Duration) {
	ex.stealInterval = interval
}
