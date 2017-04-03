package binder

import "github.com/yuin/gopher-lua"

type Push struct {
	context *Context
}

func (p *Push) String(s string) {
	p.context.state.Push(lua.LString(s))
	p.context.increase()
}

func (p *Push) Number(n float64) {
	p.context.state.Push(lua.LNumber(n))
	p.context.increase()
}

func (p *Push) Bool(b bool) {
	p.context.state.Push(lua.LBool(b))
	p.context.increase()
}

func (p *Push) Data(d interface{}, t string) {
	ud := p.context.state.NewUserData()
	ud.Value = d

	p.context.state.SetMetatable(ud, p.context.state.GetTypeMetatable(t))
	p.context.state.Push(ud)

	p.context.increase()
}
