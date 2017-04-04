package binder

import (
	"errors"
	"fmt"
	"testing"
)

const (
	testFibonacciIterations = 30
)

func BenchmarkFibonacci_Go(b *testing.B) {
	benchNative(New(Options{
		SkipOpenLibs: true,
	}), b)

}

func BenchmarkFibonacci_Go_OpenLibs(b *testing.B) {
	benchNative(New(), b)

}

func BenchmarkFibonacci_Lua(b *testing.B) {
	benchLua(New(Options{
		SkipOpenLibs: true,
	}), b)
}

func BenchmarkFibonacci_Lua_OpenLibs(b *testing.B) {
	benchLua(New(), b)
}

// fib is go native version of lua fib()
func fib(n float64) float64 {
	if n < 3 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func benchNative(bndr *Binder, b *testing.B) {
	bndr.Func("fib", func(c *Context) error {
		t := c.Top()
		if t != 1 {
			return errors.New("need an argument")
		}

		c.Push().Number(fib(c.Arg(1).Number()))
		return nil
	})

	s := `
		for i = 1, %d do
			fib(%d)
		end`

	if err := bndr.DoString(fmt.Sprintf(s, b.N, testFibonacciIterations)); err != nil {
		b.Error("Error Fibonacci go calculation", err)
	}
}

func benchLua(bndr *Binder, b *testing.B) {
	s := `
		function fib(n)
			if n<3 then
				return 1
			else
				return fib(n-1) + fib(n-2)
			end
		end

		for i = 1, %d do
			fib(%d)
		end`

	if err := bndr.DoString(fmt.Sprintf(s, b.N, testFibonacciIterations)); err != nil {
		b.Error("Error Fibonacci go calculation", err)
	}
}
