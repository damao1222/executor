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
)

type FixedExecutor struct {
	runners  []TaskRunner
	taskChan chan Task
	stop     chan bool
}

//创建一个固定大小的协程池
//size： 协程池大小
//taskBufSize： 任务缓冲大小
//需要注意的是，可以成功加入协程池的任务数量为 size + taskBufSize
//协程池会自动选择一个空闲的协程执行任务，所有任务最终都将被执行
func NewFixedExecutor(size int, taskBufSize int) *FixedExecutor {
	ex := &FixedExecutor{
		runners:  make([]TaskRunner, size),
		taskChan: make(chan Task, size+taskBufSize),
		stop:     make(chan bool),
	}
	for i := 0; i < size; i++ {
		ex.runners[i] = NewOnce()
		//start runner loop
		go ex.runners[i].Loop()
	}

	go ex.loop()

	return ex
}

func (ex *FixedExecutor) Run(task Task) error {
	select {
	case <-ex.stop:
		return errors.New("executor is stopped")
	case ex.taskChan <- task:
		return nil
	default:
		return errors.New("All Runners are busy. ")
	}
}

func (ex *FixedExecutor) loop() {
	for {
		select {
		case <-ex.stop:
			return
		case task, ok := <-ex.taskChan:
			if ok {
				for i := 0; i < len(ex.runners); i++ {
					runner := ex.runners[i]
					if runner.SetTask(task) {
						break
					}
				}
			}
		}
	}
}

func (ex *FixedExecutor) Stop() {
	//stop self loop first
	close(ex.stop)
	//close task channel
	//close(ex.taskBuf)
	//stop runner
	for _, runner := range ex.runners {
		runner.Stop()
	}
}
