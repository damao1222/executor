/**
 * Copyright (C) 2018-2020, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/16
 * @time 18:01
 * @version V1.0
 * Description:
 */

package test

import (
	"fmt"
	"github.com/xfali/executor"
	"testing"
	"time"
)

func TestLazyExecutor(t *testing.T) {
	executor := executor.NewLazyFixBUfExecutor()
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

func TestLazyExecutor2(t *testing.T) {
	executor := executor.NewLazyFixBUfExecutor()

	for i := 0; i < 100; i++ {
		b := i
		executor.Run(func() {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("func %d done \n", b)
		})
	}

	<-time.After(10 * time.Second)
}

