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
    executor := executor.NewFixedExecutor(2, 1)
    err := executor.Run(func() {
        time.Sleep(500*time.Millisecond)
        fmt.Println("func1 done")
    }, time.Second, func() {
        fmt.Println("func1 timeout")
    })

    if err != nil {
        t.Fail()
    }

    err = executor.Run(func() {
        time.Sleep(2 * time.Second)
        fmt.Println("func2 done")
    }, time.Second, func() {
        fmt.Println("func2 timeout")
    })

    if err != nil {
        t.Fail()
    }

    err = executor.Run(func() {
        time.Sleep(time.Second)
        fmt.Println("func3 done")
    }, time.Second, func() {
        fmt.Println("func2 timeout")
    })

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

func TestFixedExecutor2(t *testing.T) {
    executor := executor.NewFixedExecutor(2, 20)

    for i:=0; i<100; i++{
        b := i
        executor.Run(func() {
            time.Sleep(500*time.Millisecond)
            fmt.Printf("func %d done \n", b)
        }, time.Second, func() {
            fmt.Printf("func %d timeout \n", b)
        })
    }

    for {
        select {
        case <- time.After(10*time.Second):
            return
        }
    }
}

func TestFixedExecutor3(t *testing.T) {
    executor := executor.NewFixedExecutor(2, 20)

    executor.Run(func() {
        time.Sleep(500*time.Millisecond)
        fmt.Println("done")
    }, time.Second, nil)

    for {
        select {
        case <- time.After(10*time.Second):
            return
        }
    }
}
