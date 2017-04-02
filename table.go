package binder

import (
	"github.com/yuin/gopher-lua"
)

type Handler func(*Context) error

type Table struct {
	Name    string
	binder  *Binder
	static  map[string]Handler
	methods map[string]Handler
}

func (t *Table) Static(name string, handler Handler) {
	t.static[name] = handler
}

func (t *Table) Method(name string, handler Handler) {
	t.methods[name] = handler
}

func (t *Table) prepare(s *lua.LState) {
	mt := s.NewTypeMetatable(t.Name)
	s.SetGlobal(t.Name, mt)

	for name, handler := range t.static {
		s.SetField(mt, name, s.NewFunction(func(s *lua.LState) int {
			return t.handle(s, handler)
		}))
	}

	if len(t.methods) > 0 {
		methods := map[string]lua.LGFunction{}

		for name, handler := range t.methods {
			methods[name] = func(state *lua.LState) int {
				return t.handle(state, handler)
			}
		}

		s.SetField(mt, "__index", s.SetFuncs(s.NewTable(), methods))
	}
}

func (t *Table) handle(s *lua.LState, handler Handler) int {
	c := &Context{
		e:     true,
		state: s,
		table: t,
	}

	err := handler(c)
	if err != nil {
		c.error(err)
	}

	if !c.empty() {
		return 1
	}

	return 0
}
