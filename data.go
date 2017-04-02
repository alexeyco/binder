package binder

import "github.com/yuin/gopher-lua"

type Data struct {
	state *lua.LState
	number int
}

func (d *Data) Value() interface{} {
	return d.state.CheckUserData(d.number).Value
}
