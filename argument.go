package binder

import "github.com/yuin/gopher-lua"

type Argument struct {
	state  *lua.LState
	number int
}

func (a *Argument) String() string {
	return a.state.CheckString(a.number)
}

func (a *Argument) Number() float64 {
	return float64(a.state.CheckNumber(a.number))
}

func (a *Argument) Bool() bool {
	return a.state.CheckBool(a.number)
}

func (a *Argument) Any() interface{} {
	return a.state.CheckAny(a.number)
}

func (a *Argument) Data() interface{} {
	return a.state.CheckUserData(a.number).Value
}
