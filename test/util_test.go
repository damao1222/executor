// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/executor"
	"github.com/xfali/executor/exectool"
	"math/rand"
	"testing"
	"time"
)

func TestRunUtil(t *testing.T) {
	exec := executor.NewWorkStealingExecutor()
	t.Run("type not match", func(t *testing.T) {
		var ret []string
		_, err := exectool.Run(exec, &ret, func() int {
			return 1
		}, func() int {
			return 2
		})
		if err == nil {
			t.Fatal("must be error!")
		}
		t.Log(err)
	})

	t.Run("normal", func(t *testing.T) {
		var ret []string
		now := time.Now()
		r, err := exectool.Run(exec, &ret, func() string {
			time.Sleep(1 * time.Second)
			return "hello"
		}, func() string {
			time.Sleep(1 * time.Second)
			return "world"
		})
		if err != nil {
			t.Fatal("must not be error! ", err)
		}
		if r.HaveError() {
			t.Fatal("must not have error")
		}
		ti := time.Since(now).Milliseconds()
		if ti > int64((time.Second+(time.Millisecond*100))/time.Millisecond) {
			t.Fatal("time cost too much. ")
		}
		t.Logf("use time: %d ms", ti)
		for _, v := range ret {
			t.Log(v)
			if v != "hello" && v != "world" {
				t.Fatal("not match")
			}
		}
	})

	t.Run("panic", func(t *testing.T) {
		var ret []string
		now := time.Now()
		r, err := exectool.Run(exec, &ret, func() string {
			time.Sleep(1 * time.Second)
			return "hello"
		}, func() string {
			time.Sleep(500 * time.Millisecond)
			panic("world")
			return "world"
		})
		if err != nil {
			t.Fatal("must not be error! ", err)
		}
		if !r.HaveError() {
			t.Fatal("must have error")
		}
		if err := r.GetError(1); err == nil {
			t.Fatal("err must not nil")
		} else {
			t.Log(err)
		}
		ti := time.Since(now).Milliseconds()
		if ti > int64((time.Second+(time.Millisecond*100))/time.Millisecond) {
			t.Fatal("time cost too much. ")
		}
		t.Logf("use time: %d ms", ti)
		for _, v := range ret {
			t.Log(v)
			if v != "hello" {
				t.Fatal("not match")
			}
		}
	})

	t.Run("random", func(t *testing.T) {
		var ret []int
		now := time.Now()
		var funcs []interface{}
		for i := 0; i < executor.DefaultMaxGorutineSize; i++ {
			cur := i
			funcs = append(funcs, func() int {
				v := rand.Intn(100)
				time.Sleep(time.Millisecond * time.Duration(cur * v))
				return cur
			})
		}
		r, err := exectool.Run(exec, &ret, funcs...)
		if err != nil {
			t.Fatal("must not be error! ", err)
		}
		if r.HaveError() {
			t.Fatal("must not have error")
		}
		ti := time.Since(now).Milliseconds()
		if ti > int64((executor.DefaultMaxGorutineSize*time.Second+(time.Millisecond*100))/time.Millisecond) {
			t.Fatal("time cost too much. ")
		}
		t.Logf("use time: %d ms", ti)
		for _, v := range ret {
			t.Log(v)
		}
	})
}
