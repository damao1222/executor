// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"fmt"
	"github.com/xfali/executor"
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
