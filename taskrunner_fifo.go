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

type TaskRunnerFIFO struct {
    task chan Task
    stop chan bool
}

func NewFIFO(taskSize int) *TaskRunnerFIFO {
    return &TaskRunnerFIFO{make(chan Task, taskSize), make(chan bool)}
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
    close(tr.stop)
}

func (tr *TaskRunnerFIFO) Next() {

}

func (tr *TaskRunnerFIFO) OnExpired(task Task) {
}

func (tr *TaskRunnerFIFO) Loop() {
    for {
        select {
        case task, ok := <-tr.task:
            if ok {
                task()
            }
        case <-tr.stop:
            return
        }
    }
}
