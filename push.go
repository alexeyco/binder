package binder

import "github.com/yuin/gopher-lua"

type Push struct {
	context *Context
}

func (p *Push) Table(v interface{}) {
	if p.context.table == nil {
		p.Value(v)
		return
	}

	state := p.context.state

	ud := state.NewUserData()
	ud.Value = v

	state.SetMetatable(ud, state.GetTypeMetatable(p.context.table.Name))
	state.Push(ud)
	p.context.e = false
}

func (p *Push) Value(v interface{}) {
	var c lua.LValue

	switch t := v.(type) {
	case bool:
		c = lua.LBool(t)
	case string:
		c = lua.LString(t)
	case int:
		c = lua.LNumber(t)
	case uint:
		c = lua.LNumber(t)
	case float32, float64:
		c = lua.LNumber(t.(float64))
	default:
		c = lua.LValue(&lua.LNilType{})
	}

	p.context.state.Push(c)
	p.context.e = false
}
