package binder

import "github.com/yuin/gopher-lua"

// Push function result wrapper
type Push struct {
	context *Context
}

// String pushes sting function result
func (p *Push) String(s string) {
	p.context.state.Push(lua.LString(s))
	p.context.increase()
}

// Number pushes sting function result
func (p *Push) Number(n float64) {
	p.context.state.Push(lua.LNumber(n))
	p.context.increase()
}

// Bool pushes bool function result
func (p *Push) Bool(b bool) {
	p.context.state.Push(lua.LBool(b))
	p.context.increase()
}

// Data pushes UserData function result
func (p *Push) Data(d interface{}, t string) {
	ud := p.context.state.NewUserData()
	ud.Value = d

	p.context.state.SetMetatable(ud, p.context.state.GetTypeMetatable(t))
	p.context.state.Push(ud)

	p.context.increase()
}
