// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"bytes"
	"fmt"
	"github.com/xfali/executor"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestStealExecutor(t *testing.T) {
	executor := executor.NewWorkStealingExecutor()
	err := executor.Run(func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("func1 done")
	})

	if err != nil {
		t.Fatal(err)
	}

	err = executor.Run(func() {
		time.Sleep(2 * time.Second)
		fmt.Println("func2 done")
	})

	if err != nil {
		t.Fatal(err)
	}

	err = executor.Run(func() {
		time.Sleep(time.Second)
		fmt.Println("func3 done")
	})

	<-time.After(5 * time.Second)
}

func TestStealExecutorSteal(t *testing.T) {
	executor := executor.NewWorkStealingExecutor(
		executor.SetInitialGorutineSize(2),
		executor.SetMaxGorutineSize(2),
		executor.SetStealInterval(100*time.Millisecond))

	now := time.Now()
	err := executor.Run(func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("func1 done ", getGID(), " time ", time.Since(now) / time.Millisecond)
	})

	if err != nil {
		t.Fatal(err)
	}

	err = executor.Run(func() {
		time.Sleep(3 * time.Second)
		fmt.Println("func2 done ", getGID(), " time ", time.Since(now) / time.Millisecond)
	})

	if err != nil {
		t.Fatal(err)
	}

	err = executor.Run(func() {
		time.Sleep(time.Second)
		fmt.Println("func3 done ", getGID(), " time ", time.Since(now) / time.Millisecond)
	})

	err = executor.Run(func() {
		time.Sleep(time.Second)
		fmt.Println("func4 done ", getGID(), " time ", time.Since(now) / time.Millisecond)
	})

	<-time.After(5 * time.Second)
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
