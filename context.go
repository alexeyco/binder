package binder

import "github.com/yuin/gopher-lua"

type Context struct {
	e     bool
	state *lua.LState
	table *Table
}

func (c *Context) Param(number int) *Param {
	return &Param{
		state:  c.state,
		number: number,
	}
}

func (c *Context) Data(number int) *Data {
	return &Data{
		state:  c.state,
		number: number,
	}
}

func (c *Context) Top() int {
	return c.state.GetTop()
}

func (c *Context) Push() *Push {
	return &Push{
		context: c,
	}
}

func (c *Context) error(e error) {
	c.state.RaiseError(e.Error())
}

func (c *Context) empty() bool {
	return c.e
}
