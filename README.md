## Decorator in Golang
[![Build Status](https://travis-ci.org/starwander/degorator.svg?branch=master)](https://travis-ci.org/starwander/degorator)
[![codecov](https://codecov.io/gh/starwander/degorator/branch/master/graph/badge.svg)](https://codecov.io/gh/starwander/degorator)
[![Go Report Card](https://goreportcard.com/badge/github.com/starwander/degorator)](https://goreportcard.com/report/github.com/starwander/degorator)
[![GoDoc](https://godoc.org/github.com/starwander/degorator?status.svg)](https://godoc.org/github.com/starwander/degorator)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

Degorator implements the [decorator pattern](https://en.wikipedia.org/wiki/Decorator_pattern) in [Golang](https://golang.org/).
This can be used to add behavior, such as logs or metrics, into a function without affecting the original behavior at runtime.

##Requirements

    go get github.com/starwander/degorator

## Usage
* Decorate: injects two functions(injectedBefore & injectedAfter) into the target function.
* MakeDecorator: generate a decorator to a certain function type which can be used later.

|      Original       |          Decorated           |
| :-----------------: | :--------------------------: |
| func Myfunc(in)out{ | func MyfuncDecorated(in)out{ |
|     ......          |    injectedBefore(in)        |
|     do someting     |    out = MyFunc(in)          |
|     ......          |    injectedAfter(out)        |
| }                   | }                            |

## Example

```go
package main

import (
	"fmt"
	"github.com/starwander/degorator"
)

type Counter struct {
	number int
	error  int
}

func (m *Counter) add(s string) {
	if s == "nothing" {
		return
	}
	m.number++
}

func (m *Counter) addErr(err error) {
	if err != nil {
		m.error++
	}
}

func Log(s string) {
	fmt.Println("input:", s)
}

type MyFunc func(s string) error

var myFunc MyFunc = func(s string) error {
	if s == "error" {
		return fmt.Errorf("error")
	}
	return nil
}

func main() {
	counter := new(Counter)

	var CounterDecorator func(MyFunc) MyFunc
	err := degorator.MakeDecorator(&CounterDecorator, counter.add, counter.addErr)
	if err != nil {
		panic(err)
	}
	myFunc = CounterDecorator(myFunc)

	myFunc("something")
	//1 0
	fmt.Println(counter.number, counter.error)

	myFunc("error")
	//2 1
	fmt.Println(counter.number, counter.error)

	myFunc("nothing")
	//2 1
	fmt.Println(counter.number, counter.error)

	var myFuncDecorated MyFunc
	err = degorator.Decorate(&myFuncDecorated, myFunc, Log, nil)
	if err != nil {
		panic(err)
	}

	//input: another
	myFuncDecorated("another")
	//3 1
	fmt.Println(counter.number, counter.error)
}
```

## Reference

[GoDoc](https://godoc.org/github.com/starwander/degorator)

## LICENSE

Degorator source code is licensed under the [Apache Licence, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
