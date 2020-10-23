/**
 * Copyright (C) 2018-2020, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/16
 * @time 16:49
 * @version V1.0
 * Description:
 */

package executor

type TaskRunnerOnce struct {
	task chan Task
	stop chan bool
}

func NewOnce() *TaskRunnerOnce {
	return &TaskRunnerOnce{make(chan Task), make(chan bool)}
}

//NOTICE:当Loop协程没有就绪，则会一直返回false
func (tr *TaskRunnerOnce) SetTask(task Task) bool {
	select {
	case tr.task <- task:
		return true
	default:
		return false
	}
}

func (tr *TaskRunnerOnce) Stop() {
	close(tr.stop)
}

func (tr *TaskRunnerOnce) Next() {

}

func (tr *TaskRunnerOnce) OnExpired(Task) {

}

func (tr *TaskRunnerOnce) Loop() {
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
