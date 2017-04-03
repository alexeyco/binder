package binder

import (
	"github.com/yuin/gopher-lua"
)

// Module is a lua module wrapper
type Module struct {
	name   string
	state  *lua.LState
	fields map[string]lua.LValue
	funcs  map[string]Handler
}

// String sets module string constant
func (m *Module) String(name, value string) {
	m.fields[name] = lua.LString(value)
}

// Number sets module number (float64) constant
func (m *Module) Number(name string, value float64) {
	m.fields[name] = lua.LNumber(value)
}

// Bool sets module bool constant
func (m *Module) Bool(name string, value bool) {
	m.fields[name] = lua.LBool(value)
}

// Func sets module function with specified name
func (m *Module) Func(name string, handler Handler) {
	m.funcs[name] = handler
}

func (m *Module) load() {
	m.state.PreloadModule(m.name, func(state *lua.LState) int {
		module := state.SetFuncs(state.NewTable(), exports(m.funcs))

		for name, value := range m.fields {
			state.SetField(module, name, value)
		}

		state.Push(module)
		return 1
	})
}
