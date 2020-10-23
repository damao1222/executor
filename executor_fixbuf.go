/**
 * Copyright (C) 2018-2020, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/16
 * @time 16:35
 * @version V1.0
 * Description:
 */

package executor

import (
	"errors"
	"sync"
)

type FixedBufExecutor struct {
	runners  []TaskRunner
	curIndex int
	lock     sync.Mutex
}

//创建一个固定大小的协程池
//size： 协程池大小
//taskBufSize： 任务缓冲大小
//需要注意的是，可以成功加入协程池的任务数量为 size * taskBufSize
//并且当第一个协程池中的协程taskBufSize被填满之前，不会尝试使用协程池中的其他协程
func NewFixedBufExecutor(size int, taskBufSize int) *FixedBufExecutor {
	ex := &FixedBufExecutor{
		runners:  make([]TaskRunner, size),
		curIndex: 0,
	}
	for i := 0; i < size; i++ {
		ex.runners[i] = NewFIFO(taskBufSize)
		//start runner loop
		go ex.runners[i].Loop()
	}
	return ex
}

func (ex *FixedBufExecutor) Run(task Task) error {
	for i := 0; i < len(ex.runners); i++ {
		//runner := ex.runners[i]
		runner := ex.selectRunner()
		if runner.SetTask(task) {
			return nil
		}
	}
	return errors.New("All Runners are busy. ")
}

func (ex *FixedBufExecutor) selectRunner() TaskRunner {
	ex.lock.Lock()
	defer ex.lock.Unlock()

	ex.curIndex++
	ex.curIndex = ex.curIndex % len(ex.runners)

	return ex.runners[ex.curIndex]
}

func (ex *FixedBufExecutor) Stop() {
	for _, runner := range ex.runners {
		runner.Stop()
	}
}
