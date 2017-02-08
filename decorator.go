// Copyright(c) 2017 Ethan Zhuang <zhuangwj@gmail.com>.

package degorator

import (
	"fmt"
	"reflect"
)

func Decorate(decorated interface{}, target interface{}, before interface{}, after interface{}) (err error) {
	var targetFunc reflect.Value
	var decoratedFunc reflect.Value
	var beforeFunc reflect.Value
	var afterFunc reflect.Value

	decoratedFunc, err = checkFPTR(decorated)
	if err != nil {
		return
	}

	targetFunc = reflect.ValueOf(target)
	if targetFunc.Kind() != reflect.Func {
		err = fmt.Errorf("Input target para is not a function.")
		return
	}

	beforeFunc, afterFunc, err = checkInjection(targetFunc.Type(), before, after)
	if err != nil {
		return
	}

	decoratedFunc.Set(reflect.MakeFunc(targetFunc.Type(), func(in []reflect.Value) (out []reflect.Value) {
		if targetFunc.Type().IsVariadic() {
			if before != nil {
				beforeFunc.CallSlice(in)
			}
			out = targetFunc.CallSlice(in)
		} else {
			if before != nil {
				beforeFunc.Call(in)
			}
			out = targetFunc.Call(in)
		}
		if after != nil {
			afterFunc.Call(out)
		}
		return
	}))
	return
}

func MakeDecorator(decorator interface{}, before interface{}, after interface{}) (err error) {
	var decoFunc reflect.Value
	var beforeFunc reflect.Value
	var afterFunc reflect.Value

	decoFunc, err = checkDecorator(decorator)
	if err != nil {
		return
	}

	beforeFunc, afterFunc, err = checkInjection(decoFunc.Type().In(0), before, after)
	if err != nil {
		return
	}

	decoFunc.Set(reflect.MakeFunc(decoFunc.Type(), func(args []reflect.Value) (results []reflect.Value) {
		wrappedFunc := func(in []reflect.Value) (out []reflect.Value) {
			if args[0].Type().IsVariadic() {
				if before != nil {
					beforeFunc.CallSlice(in)
				}
				out = args[0].CallSlice(in)
			} else {
				if before != nil {
					beforeFunc.Call(in)
				}
				out = args[0].Call(in)
			}
			if after != nil {
				afterFunc.Call(out)
			}
			return
		}
		v := reflect.MakeFunc(args[0].Type(), wrappedFunc)
		results = []reflect.Value{v}
		return
	}))
	return
}

func checkFPTR(fptr interface{}) (function reflect.Value, err error) {
	if fptr == nil {
		err = fmt.Errorf("Input para is nil.")
		return
	}
	if reflect.TypeOf(fptr).Kind() != reflect.Ptr {
		err = fmt.Errorf("Input para is not a pointer.")
		return
	}
	function = reflect.ValueOf(fptr).Elem()
	if function.Kind() != reflect.Func {
		err = fmt.Errorf("Input para is not a pointer to a function.")
		return
	}
	return
}

func checkInjection(targetType reflect.Type, before interface{}, after interface{}) (beforeFunc reflect.Value, afterFunc reflect.Value, err error) {
	if before != nil {
		beforeFunc = reflect.ValueOf(before)
		if beforeFunc.Kind() != reflect.Func {
			err = fmt.Errorf("Only a function can be injected before.")
			return
		}
		if beforeFunc.Type().NumIn() != targetType.NumIn() {
			err = fmt.Errorf("The input para number of the funtion injected before must be same with the input para number of the target funtion.")
			return
		}
		for i := 0; i < beforeFunc.Type().NumIn(); i++ {
			if beforeFunc.Type().In(i) != targetType.In(i) {
				err = fmt.Errorf("The input para types of the funtion injected before must be same with the input para types of the target funtion.")
				return
			}
		}
	}
	if after != nil {
		afterFunc = reflect.ValueOf(after)
		if afterFunc.Kind() != reflect.Func {
			err = fmt.Errorf("Only a function can be injected after.")
			return
		}
		if afterFunc.Type().NumIn() != targetType.NumOut() {
			err = fmt.Errorf("The input para number of the funtion injected after must be same with the output para number of the target funtion.")
			return
		}
		for i := 0; i < afterFunc.Type().NumIn(); i++ {
			if afterFunc.Type().In(i) != targetType.Out(i) {
				err = fmt.Errorf("The input para types of the funtion injected after must be same with the output para types of the target funtion.")
				return
			}
		}
	}
	return
}

func checkDecorator(decorator interface{}) (decoFunc reflect.Value, err error) {
	decoFunc, err = checkFPTR(decorator)
	if err != nil {
		return
	}
	if decoFunc.Type().NumIn() != 1 || decoFunc.Type().NumOut() != 1 {
		err = fmt.Errorf("Decorator function must have one input para and one output para.")
		return
	}
	if decoFunc.Type().In(0).Kind() != reflect.Func || decoFunc.Type().Out(0).Kind() != reflect.Func {
		err = fmt.Errorf("Decorator function's input para type and output para type must be function type.")
		return
	}
	if decoFunc.Type().In(0).NumIn() != decoFunc.Type().NumOut() {
		err = fmt.Errorf("Decoratee function must have same input and ouput para number.")
		return
	}
	for i := 0; i < decoFunc.Type().In(0).NumIn(); i++ {
		if decoFunc.Type().In(0).In(i) != decoFunc.Type().Out(0).In(i) {
			err = fmt.Errorf("Decoratee function must have same input and output para type.")
			return
		}
	}
	return
}
