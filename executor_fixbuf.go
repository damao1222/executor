/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/16 
 * @time 16:35
 * @version V1.0
 * Description: 
 */

package executor

import (
    "time"
    "github.com/xfali/timewheel"
    "github.com/xfali/timewheel/async"
    "errors"
    "sync/atomic"
)


type FixedBufExecutor struct {
    timewheel timewheel.TimeWheel
    runners []TaskRunner
    curIndex uint32
}

//创建一个固定大小的协程池
//size： 协程池大小
//taskBufSize： 任务缓冲大小
//需要注意的是，可以成功加入协程池的任务数量为 size * taskBufSize
//并且当第一个协程池中的协程taskBufSize被填满之前，不会尝试使用协程池中的其他协程
func NewFixedBufExecutor(size int, taskBufSize int) Executor {
    ex := &FixedBufExecutor{
        timewheel : async.New(20*time.Millisecond, time.Minute),
        runners : make([]TaskRunner, size),
        curIndex: 0,
    }
    //start timer
    ex.timewheel.Start()
    for i:=0; i<size; i++ {
        ex.runners[i] = NewFIFO(taskBufSize)
        //start runner loop
        go ex.runners[i].Loop()
    }
    return ex
}

func (ex *FixedBufExecutor)Run(task Task, expire time.Duration, timeout TaskTimeout) error {
    for i:=0; i<len(ex.runners); i++ {
        //runner := ex.runners[i]
        runner := ex.selectRunner()
        if runner.SetTask(task) {
            if timeout != nil {
                ex.timewheel.Add(timewheel.NewTimer(func(data interface{}) {
                    timeout()
                }, expire, nil))
            }
            return nil
        }
    }
    return errors.New("All Runners are busy")
}

//Not thread-safe but that's OK
func (ex *FixedBufExecutor)selectRunner() TaskRunner {
    index := atomic.AddUint32(&ex.curIndex, 1) % uint32(len(ex.runners))
    //由于在多协程环境，某一个task runner是否idle是动态变化的，所以没必要保证此处的线程安全
    //ex.curIndex++
    //index := int(ex.curIndex) % len(ex.runners)
    return ex.runners[index]
}

func (ex *FixedBufExecutor)Stop() {
    for _, runner := range ex.runners {
        runner.Stop()
    }
}
