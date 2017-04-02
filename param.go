package binder

import "github.com/yuin/gopher-lua"

type Param struct {
	state  *lua.LState
	number int
}

func (p *Param) String() string {
	return p.state.CheckString(p.number)
}

func (p *Param) Int() int {
	return p.state.CheckInt(p.number)
}

func (p *Param) Bool() bool {
	return p.state.CheckBool(p.number)
}

func (p *Param) Any() interface{} {
	return p.state.CheckAny(p.number)
}
