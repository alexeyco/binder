package binder

import (
	"errors"
	"fmt"
	"testing"
)

const (
	testFibonacciIterations = 30
)

func TestBenchmarks_Assert(t *testing.T) {
	bndr := New()

	bndr.Func("fib_go", func(c *Context) error {
		t := c.Top()
		if t != 1 {
			return errors.New("need an argument")
		}

		c.Push().Number(fib(c.Arg(1).Number()))
		return nil
	})

	if err := bndr.DoString(`
		function fib_lua(n)
			if n<3 then
				return 1
			else
				return fib_lua(n-1) + fib_lua(n-2)
			end
		end

		assert(fib_go(3) == fib_lua(3), 'fib_go(3) != fib_lua(3)')
		assert(fib_go(10) == fib_lua(10), 'fib_go(10) != fib_lua(10)')
		assert(fib_go(30) == fib_lua(30), 'fib_go(30) != fib_lua(30)')
	`); err != nil {
		t.Error("Lua and go funcs are not equal", err)
	}
}

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
	}

	return fib(n-1) + fib(n-2)
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
