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

type Task func()
type TaskTimeout func()

type Executor interface {
	//执行一个任务
	//Task：在某一个协程中执行的任务
	Run(Task) error

	//停止协程池
	Stop()
}
