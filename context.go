package binder

import "github.com/yuin/gopher-lua"

// Context function context
type Context struct {
	state  *lua.LState
	pushed int
}

// Top returns count of function arguments
func (c *Context) Top() int {
	return c.state.GetTop()
}

// Arg returns function argument by number
func (c *Context) Arg(num int) *Argument {
	return &Argument{
		state:  c.state,
		number: num,
	}
}

// Push pushes function result
func (c *Context) Push() *Push {
	return &Push{
		context: c,
	}
}

func (c *Context) increase() {
	c.pushed++
}

func (c *Context) error(e string) {
	c.state.RaiseError(e)
}
