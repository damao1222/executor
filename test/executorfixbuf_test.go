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
    "github.com/xfali/executor"
    "time"
    "fmt"
    "github.com/xfali/goutils/atomic"
    "math"
    atomic2 "sync/atomic"
)

func TestFixedBufExecutor(t *testing.T) {
    executor := executor.NewFixedBufExecutor(2, 1)
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

    select {
    case <- time.After(5*time.Second):
        return
    }
}

func TestFixedBufExecutor2(t *testing.T) {
    executor := executor.NewFixedBufExecutor(2, 20)

    for i:=0; i<100; i++{
        b := i
        executor.Run(func() {
            time.Sleep(500*time.Millisecond)
            fmt.Printf("func %d done \n", b)
        }, time.Second, func() {
            fmt.Printf("func %d timeout \n", b)
        })
    }

    select {
    case <- time.After(10*time.Second):
        return
    }
}

func TestFixedBufExecutor3(t *testing.T) {
    executor := executor.NewFixedBufExecutor(2, 20)

    executor.Run(func() {
        time.Sleep(500*time.Millisecond)
        fmt.Println("done")
    }, time.Second, nil)

    select {
    case <- time.After(10*time.Second):
        return
    }
}

func TestFixedBufExecutor4(t *testing.T) {
    executor := executor.NewFixedBufExecutor(100, 20)

    count := 0
    for i:=1; i<100; i++ {
        b := i
        e := atomic.AtomicBool(0)
        now := time.Now()
        if executor.Run(func() {
            fmt.Printf("%d set %t\n", b, e.IsSet())
            if e.IsSet() {
                fmt.Printf("%d expired\n", b)
                return
            }
            fmt.Printf("%d done time %d\n", b, time.Since(now) / time.Millisecond)
            time.Sleep(time.Duration(b*100)*time.Millisecond)
        }, time.Second, func() {
            fmt.Printf("call expire %d at: %d\n",b, time.Since(now) / time.Millisecond)
            e.Set()
        }) == nil {
            count++
        } else {
            fmt.Printf("run error")
        }
    }

    select {
    case <- time.After(10*time.Second):
        return
    }
    fmt.Printf("success count %d\n", count)
}

func TestFixedBufExecutor5(t *testing.T) {
    executor := executor.NewFixedBufExecutor(100, 0)
    //wait for runner ready
    time.Sleep(time.Second)

    count := 0
    for i:=1; i<100; i++ {
        b := i
        e := atomic.AtomicBool(0)
        now := time.Now()
        if executor.Run(func() {
            fmt.Printf("%d set %t\n", b, e.IsSet())
            if e.IsSet() {
                fmt.Printf("%d expired\n", b)
                return
            }
            fmt.Printf("%d done time %d\n", b, time.Since(now) / time.Millisecond)
            time.Sleep(time.Duration(b*100)*time.Millisecond)
        }, time.Second, func() {
            fmt.Printf("call expire %d at: %d\n",b, time.Since(now) / time.Millisecond)
            e.Set()
        }) == nil {
            count++
        } else {
            fmt.Printf("run error\n")
        }
    }

    select {
    case <- time.After(10*time.Second):
        return
    }
    fmt.Printf("success count %d\n", count)
}

func TestFixedBufExecutor6(t *testing.T) {
    executor := executor.NewFixedBufExecutor(50, 20)

    count := 0
    for i:=1; i<100; i++ {
        b := i
        e := atomic.AtomicBool(0)
        now := time.Now()
        if executor.Run(func() {
            fmt.Printf("%d set %t\n", b, e.IsSet())
            if e.IsSet() {
                fmt.Printf("%d expired\n", b)
                return
            }
            fmt.Printf("%d done time %d\n", b, time.Since(now) / time.Millisecond)
            time.Sleep(time.Duration(500)*time.Millisecond)
        }, time.Second, func() {
            fmt.Printf("call expire %d at: %d\n",b, time.Since(now) / time.Millisecond)
            e.Set()
        }) == nil {
            count++
        } else {
            fmt.Printf("run error")
        }
    }

    select {
    case <- time.After(10*time.Second):
        return
    }
    fmt.Printf("success count %d\n", count)
}

func TestUint(t *testing.T) {
    var i uint32
    b := false
    for {
        //i++
        atomic2.AddUint32(&i, 1)
        if b {
            fmt.Println(i)
            return
        }
        if i == math.MaxUint32 {
            fmt.Println("MAX", i)
            b = true
        }
    }
}