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

type TaskState int32

const (
	TaskStateIdle    = 0
	TaskStateRunning = 1
)

type TaskRunner interface {
	// 尝试添加一个任务
	// Task：任务
	// 成功返回true，失败返回false
	SetTask(Task) bool

	// 是否有任务正在执行
	IsIdle() bool

	// 停止任务执行器
	Stop()

	// 执行器循环
	Loop()

	// 设置任务状态通知器
	SetNotifier(chan<- TaskRunner)
}
