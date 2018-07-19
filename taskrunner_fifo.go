/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/16 
 * @time 16:49
 * @version V1.0
 * Description: 
 */

package executor

import (
    "github.com/damao/timewheel/utils"
)

type TaskRunnerFIFO struct {
    task chan Task
    stop utils.AtomicBool
}

func NewFIFO(taskSize int) *TaskRunnerFIFO {
    return &TaskRunnerFIFO{ make(chan Task, taskSize), 0}
}

func (tr *TaskRunnerFIFO) SetTask(task Task) bool {
    select {
    case tr.task <- task:
        return true
    default:
        return false
    }
}

func (tr *TaskRunnerFIFO) Stop() {
    tr.stop.Set()
}

func (tr *TaskRunnerFIFO) Next() {
    
}

func (tr *TaskRunnerFIFO) OnExpired(task Task) {
}

func (tr *TaskRunnerFIFO)Loop() {
    for {
        if tr.stop.IsSet() {
            break
        }

        select {
        case task, ok := <- tr.task:
            if ok && !tr.stop.IsSet() {
                task()
            }
        }
    }
}
