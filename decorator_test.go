// Copyright(c) 2017 Ethan Zhuang <zhuangwj@gmail.com>.

package degorator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"errors"
)

var _ = Describe("Tests of Decorate api", func() {
	type MyFunc func(string) (error)
	type MyFuncSlice func(...string) (error)
	var (
		myCounter *MyCounter

		myFunc MyFunc = func(s string) error {
			if s == "error" {
				return errors.New("error")
			}
			return nil
		}
		myFuncDecorated MyFunc

		myFuncSlice MyFuncSlice = func(s ...string) error {
			if s[0] == "error" {
				return errors.New("error")
			}
			return nil
		}
		myFuncSliceDecorated MyFuncSlice
	)

	Context("Tests of inputs check", func() {
		BeforeEach(func() {
		})

		AfterEach(func() {
		})

		It("Given an target function, when inject functions with different input or output para number, then return error", func() {
			Expect(Decorate(&myFuncDecorated, myFunc, nil, nil)).ShouldNot(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, nil, nil, nil)).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, "func", nil, nil)).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func() {}, nil)).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func(s1 string, s2 string) {}, nil)).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, nil, func() {})).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, nil, func(e1 error, e2 error) {})).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func(s string) {}, func(e1 error, e2 error) {})).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func(s1 string, s2 string) {}, func(e error) {})).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func(s1 string, s2 string) {}, func(e1 error, e2 error) {})).Should(HaveOccurred())
		})

		It("Given an target function, when inject functions with different input or output para type, then return error", func() {
			Expect(Decorate(&myFuncDecorated, myFunc, func(s string) {}, func(e error) {})).ShouldNot(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func(e error) {}, func(e error) {})).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func(s string) {}, func(s string) {})).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func(e error) {}, func(s string) {})).Should(HaveOccurred())
			Expect(Decorate(&myFuncDecorated, myFunc, func(s ...string) {}, nil)).Should(HaveOccurred())
		})
	})
	Context("Tests of wrapped functions", func() {
		BeforeEach(func() {
			myCounter = new(MyCounter)
		})

		AfterEach(func() {
		})

		It("Given an target function, when inject function before, then the injected function should be invoked before target function", func() {
			Expect(Decorate(&myFuncDecorated, myFunc, myCounter.addCount, nil)).ShouldNot(HaveOccurred())
			myFuncDecorated("1")
			Expect(myCounter.number).Should(Equal(1))
			myFuncDecorated("11")
			myFuncDecorated("111")
			Expect(myCounter.number).Should(Equal(3))
			myFuncDecorated("nothing")
			Expect(myCounter.number).Should(Equal(3))
		})

		It("Given an target function with variable args, when inject function before, then the injected function should be invoked before target function", func() {
			Expect(Decorate(&myFuncSliceDecorated, myFuncSlice, myCounter.addCountSlice, nil)).ShouldNot(HaveOccurred())
			myFuncSliceDecorated("1")
			Expect(myCounter.number).Should(Equal(1))
			myFuncSliceDecorated("11")
			myFuncSliceDecorated("111")
			Expect(myCounter.number).Should(Equal(3))
			myFuncSliceDecorated("nothing", "happened")
			Expect(myCounter.number).Should(Equal(3))
		})

		It("Given an target function, when inject function after, then the injected function should be invoked after target function", func() {
			Expect(Decorate(&myFuncDecorated, myFunc, nil, myCounter.addError)).ShouldNot(HaveOccurred())
			myFuncDecorated("1")
			Expect(myCounter.error).Should(Equal(0))
			myFuncDecorated("error")
			Expect(myCounter.error).Should(Equal(1))
			myFuncDecorated("11")
			Expect(myCounter.number).Should(Equal(0))
		})

		It("Given an target function, when inject functions before and after, then the injected functions should be invoked correctly", func() {
			Expect(Decorate(&myFuncSliceDecorated, myFuncSlice, myCounter.addCountSlice, myCounter.addError)).ShouldNot(HaveOccurred())
			myFuncSliceDecorated("1")
			Expect(myCounter.error).Should(Equal(0))
			myFuncSliceDecorated("error")
			Expect(myCounter.error).Should(Equal(1))
			myFuncSliceDecorated("11")
			Expect(myCounter.error).Should(Equal(1))
			myFuncSliceDecorated("nothing")
			Expect(myCounter.number).Should(Equal(3))
			Expect(myCounter.error).Should(Equal(1))
		})
	})
})

var _ = Describe("Tests of MakeDecorator api", func() {
	type MyFunc func(string) (error)
	type MyFuncSlice func(...string) (error)
	var (
		myCounter *MyCounter
	        myDecorator func(MyFunc) MyFunc
	        mySliceDecorator func(MyFuncSlice) MyFuncSlice

		myFunc MyFunc = func(s string) error {
			if s == "error" {
				return errors.New("error")
			}
			return nil
		}

		myFuncSlice MyFuncSlice = func(s ...string) error {
			if s[0] == "error" {
				return errors.New("error")
			}
			return nil
		}
	)

	Context("Tests of inputs check", func() {
		BeforeEach(func() {
		})

		AfterEach(func() {
		})

		It("Given an target function, when inject functions with different input or output para number, then return error", func() {
			Expect(MakeDecorator(&myDecorator, nil, nil)).ShouldNot(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func() {}, nil)).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func(s1 string, s2 string) {}, nil)).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, nil, func() {})).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, nil, func(e1 error, e2 error) {})).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func(s string) {}, func(e1 error, e2 error) {})).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func(s1 string, s2 string) {}, func(e error) {})).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func(s1 string, s2 string) {}, func(e1 error, e2 error) {})).Should(HaveOccurred())
		})

		It("Given an target function, when inject functions with different input or output para type, then return error", func() {
			Expect(MakeDecorator(&myDecorator, func(s string) {}, func(e error) {})).ShouldNot(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func(e error) {}, func(e error) {})).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func(s string) {}, func(s string) {})).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func(e error) {}, func(s string) {})).Should(HaveOccurred())
			Expect(MakeDecorator(&myDecorator, func(s ...string) {}, nil)).Should(HaveOccurred())
		})
	})
	Context("Tests of wrapped functions", func() {
		BeforeEach(func() {
			myCounter = new(MyCounter)
		})

		AfterEach(func() {
		})

		It("Given an target function, when inject function before, then the injected function should be invoked before target function", func() {
			Expect(MakeDecorator(&myDecorator, myCounter.addCount, nil)).ShouldNot(HaveOccurred())
			myFunc = myDecorator(myFunc)
			myFunc("1")
			Expect(myCounter.number).Should(Equal(1))
			myFunc("11")
			myFunc("111")
			Expect(myCounter.number).Should(Equal(3))
			myFunc("nothing")
			Expect(myCounter.number).Should(Equal(3))
		})

		It("Given an target function with variable args, when inject function before, then the injected function should be invoked before target function", func() {
			Expect(MakeDecorator(&mySliceDecorator, myCounter.addCountSlice, nil)).ShouldNot(HaveOccurred())
			myFuncSlice = mySliceDecorator(myFuncSlice)
			myFuncSlice("1")
			Expect(myCounter.number).Should(Equal(1))
			myFuncSlice("11")
			myFuncSlice("111")
			Expect(myCounter.number).Should(Equal(3))
			myFuncSlice("nothing", "happened")
			Expect(myCounter.number).Should(Equal(3))
		})

		It("Given an target function, when inject function after, then the injected function should be invoked after target function", func() {
			Expect(MakeDecorator(&myDecorator, nil, myCounter.addError)).ShouldNot(HaveOccurred())
			myFunc = myDecorator(myFunc)
			myFunc("1")
			Expect(myCounter.error).Should(Equal(0))
			myFunc("error")
			Expect(myCounter.error).Should(Equal(1))
			myFunc("11")
			Expect(myCounter.number).Should(Equal(0))
		})

		It("Given an target function, when inject functions before and after, then the injected functions should be invoked correctly", func() {
			Expect(MakeDecorator(&mySliceDecorator, myCounter.addCountSlice, myCounter.addError)).ShouldNot(HaveOccurred())
			myFuncSlice = mySliceDecorator(myFuncSlice)
			myFuncSlice("1")
			Expect(myCounter.error).Should(Equal(0))
			myFuncSlice("error")
			Expect(myCounter.error).Should(Equal(1))
			myFuncSlice("11")
			Expect(myCounter.error).Should(Equal(1))
			myFuncSlice("nothing")
			Expect(myCounter.number).Should(Equal(3))
			Expect(myCounter.error).Should(Equal(1))
		})
	})
})

type MyCounter struct {
	number int
	error  int
}

func (m *MyCounter) addCount(s string) {
	if s == "nothing" {
		return
	}
	m.number++
}

func (m *MyCounter) addError(err error) {
	if err != nil {
		m.error++
	}
}

func (m *MyCounter) addCountSlice(s ...string) {
	if s[0] == "nothing" {
		return
	}
	m.number++
}
