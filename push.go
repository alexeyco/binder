package binder

import "github.com/yuin/gopher-lua"

type Push struct {
	context *Context
}

func (p *Push) String(s string) {
	p.context.state.Push(lua.LString(s))
	p.context.e = false
}

func (p *Push) Data(v interface{}) {
	state := p.context.state

	ud := state.NewUserData()
	ud.Value = v

	state.SetMetatable(ud, state.GetTypeMetatable(p.context.table.Name))
	state.Push(ud)
	p.context.e = false
}
