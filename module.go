package binder

import (
	"github.com/yuin/gopher-lua"
)

type Module struct {
	Name  string
	funcs map[string]Handler
}

func (m *Module) Func(name string, handler Handler) {
	m.funcs[name] = handler
}

func (m *Module) load(s *lua.LState) int {
	f := map[string]lua.LGFunction{}

	for name, handler := range m.funcs {
		f[name] = func(state *lua.LState) int {
			return m.handle(state, handler)
		}
	}

	mod := s.SetFuncs(s.NewTable(), f)
	s.Push(mod)

	return 1
}

func (m *Module) prepare(s *lua.LState) {
	s.PreloadModule(m.Name, m.load)
}

func (m *Module) handle(s *lua.LState, handler Handler) int {
	c := &Context{
		e:     true,
		state: s,
	}

	err := handler(c)
	if err != nil {
		c.error(err)
		return 0
	}

	if !c.empty() {
		return 1
	}

	return 0
}
