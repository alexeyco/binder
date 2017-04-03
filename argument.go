package binder

import "github.com/yuin/gopher-lua"

// Argument is a call function argument
type Argument struct {
	state  *lua.LState
	number int
}

// String checks if function argument is string and return it
func (a *Argument) String() string {
	return a.state.CheckString(a.number)
}

// Number checks if function argument is number (float64) and return it
func (a *Argument) Number() float64 {
	return float64(a.state.CheckNumber(a.number))
}

// Bool checks if function argument is bool and return it
func (a *Argument) Bool() bool {
	return a.state.CheckBool(a.number)
}

// Any returns function argument as interface{}
func (a *Argument) Any() interface{} {
	return a.state.CheckAny(a.number)
}

// Data checks if function argument is UserData and return it
func (a *Argument) Data() interface{} {
	return a.state.CheckUserData(a.number).Value
}
