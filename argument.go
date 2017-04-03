package binder

import "github.com/yuin/gopher-lua"

type Argument struct {
	state  *lua.LState
	number int
}

func (a *Argument) String(n int) string {
	return a.state.CheckString(n)
}

func (a *Argument) Number(n int) float64 {
	return float64(a.state.CheckNumber(n))
}

func (a *Argument) Bool(n int) bool {
	return a.state.CheckBool(n)
}

func (a *Argument) Data(n int) interface{} {
	return a.state.CheckUserData(n)
}
