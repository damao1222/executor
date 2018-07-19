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

type TaskRunner interface {
    SetTask(Task) (bool)
    Stop()
    Loop()
}

