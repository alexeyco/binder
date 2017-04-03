package binder

import "github.com/yuin/gopher-lua"

type Context struct {
	state  *lua.LState
	pushed int
}

func (c *Context) Top() int {
	return c.state.GetTop()
}

func (c *Context) Arg(num int) *Argument {
	return &Argument{
		state:  c.state,
		number: num,
	}
}

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
