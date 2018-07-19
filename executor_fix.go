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
)


type FixedExecutor struct {
    timewheel timewheel.TimeWheel
    runners []TaskRunner
}

func NewFixedExecutor(size int, taskBufSize int) Executor {
    ex := &FixedExecutor{
        timewheel : async.New(20*time.Millisecond, time.Minute),
        runners : make([]TaskRunner, size),
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

func (ex *FixedExecutor)Run(task Task, expire time.Duration, timeout TaskTimeout) error {
    for i:=0; i<len(ex.runners); i++ {
        runner := ex.runners[i]
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

func (ex *FixedExecutor)Stop() {
    for _, runner := range ex.runners {
        runner.Stop()
    }
}
