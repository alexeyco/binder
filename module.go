package binder

import (
	"github.com/yuin/gopher-lua"
)

type Module struct {
	name   string
	state  *lua.LState
	fields map[string]lua.LValue
	funcs  map[string]Handler
}

func (m *Module) String(name, value string) {
	m.fields[name] = lua.LString(value)
}

func (m *Module) Number(name string, value float64) {
	m.fields[name] = lua.LNumber(value)
}

func (m *Module) Bool(name string, value bool) {
	m.fields[name] = lua.LBool(value)
}

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
