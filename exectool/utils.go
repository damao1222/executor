// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package exectool

import (
	"errors"
	"fmt"
	"github.com/xfali/executor"
	"reflect"
	"sync"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()

type Result []error

func RunAll(exec executor.Executor, ret interface{}, supplyFuncs ...interface{}) (Result, error) {
	if exec == nil {
		return nil, errors.New("Executor is nil. ")
	}
	t, err := checkSlice(ret)
	if err != nil {
		return nil, err
	}
	for _, f := range supplyFuncs {
		err = checkFunc(f, t)
		if err != nil {
			return nil, err
		}
	}

	errs := make([]error, len(supplyFuncs))

	v := reflect.ValueOf(ret)
	v = v.Elem()
	runner := runner{
		v: v,
	}
	wait := sync.WaitGroup{}
	wait.Add(len(supplyFuncs))
	for i := range supplyFuncs {
		cur := i
		exec.Run(func() {
			defer wait.Done()
			defer func() {
				if o := recover(); o != nil {
					errs[cur] = fmt.Errorf("Exec error: %v ", o)
				}
			}()
			runner.appendValue(runFunc(supplyFuncs[cur]))
		})
	}
	wait.Wait()

	runner.setSliceValue(v)
	return errs[:], nil
}

func (r Result) HaveError() bool {
	for _, err := range r {
		if err != nil {
			return true
		}
	}
	return false
}

func (r Result) GetError(index int) error {
	return r[index]
}

type runner struct {
	v    reflect.Value
	lock sync.Mutex
}

func (r *runner) appendValue(v reflect.Value) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.v = reflect.Append(r.v, v)
}

func (r *runner) setSliceValue(v reflect.Value) {
	r.lock.Lock()
	defer r.lock.Unlock()

	v.Set(r.v)
}

func checkSlice(o interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Ptr {
		return t, errors.New("Must be slice ptr. ")
	}
	t = t.Elem()
	if t.Kind() == reflect.Slice {
		return t.Elem(), nil
	}
	return t, errors.New("Must be slice ptr. ")
}

func checkFunc(o interface{}, elemType reflect.Type) error {
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Func {
		return errors.New("Must be function. ")
	}
	if t.NumOut() == 1 {
		if !t.Out(0).AssignableTo(elemType) {
			return fmt.Errorf("Type %s is not assignable to %s . ", t.Out(0), elemType)
		}
		return nil
	}
	return errors.New("Not support function type. ")
}

func runFunc(o interface{}) reflect.Value {
	v := reflect.ValueOf(o)
	return v.Call(nil)[0]
}
