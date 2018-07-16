/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/16 
 * @time 18:01
 * @version V1.0
 * Description: 
 */

package test

import (
    "testing"
    "github.com/damao/executor"
    "time"
    "fmt"
)

func TestFixedExecutor(t *testing.T) {
    executor := executor.NewFixedExecutor(2)
    err := executor.Run(func() {
        time.Sleep(time.Second)
        fmt.Println("func1 done")
    }, time.Second)

    if err != nil {
        t.Fail()
    }

    err = executor.Run(func() {
        time.Sleep(time.Second)
        fmt.Println("func2 done")
    }, time.Second)

    if err != nil {
        t.Fail()
    }

    err = executor.Run(func() {
        time.Sleep(time.Second)
        fmt.Println("func3 done")
    }, time.Second)

    if err == nil {
        t.Fail()
    }

    for {
        select {
        case <- time.After(5*time.Second):
            return
        }
    }
}
