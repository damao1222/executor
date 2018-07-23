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
)

type Task func()
type TaskTimeout func()

type Executor interface {
    //执行一个任务
    //Task：在某一个协程中执行的任务
    //Duration: 任务超时时间
    //TaskTimeout：从调用Run函数开始，在Duration时间之后回调，如果为nil则不回调
    Run(Task, time.Duration, TaskTimeout) error

    //停止协程池
    Stop()
}

