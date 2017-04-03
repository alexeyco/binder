package binder

import (
	"github.com/yuin/gopher-lua"
)

type Handler func(*Context) error

type Binder struct {
	state   *lua.LState
	funcs   map[string]Handler
	modules []*Module
}

func (b *Binder) Func(name string, handler Handler) {
	b.funcs[name] = handler
}

func (b *Binder) Module(name string) *Module {
	m := &Module{
		name:   name,
		state:  b.state,
		fields: map[string]lua.LValue{},
		funcs:  map[string]Handler{},
	}

	b.modules = append(b.modules, m)
	return m
}

func (b *Binder) ExecString(s string) error {
	b.load()
	return b.state.DoString(s)
}

func (b *Binder) ExecFile(f string) error {
	b.load()
	return b.state.DoFile(f)
}

func (b *Binder) load() {
	funcs := exports(b.funcs)
	for name, f := range funcs {
		b.state.SetGlobal(name, b.state.NewFunction(f))
	}

	for _, m := range b.modules {
		m.load()
	}
}

func New() *Binder {
	s := lua.NewState()
	defer s.Close()

	return &Binder{
		state:   s,
		funcs:   map[string]Handler{},
		modules: []*Module{},
	}
}

func exports(funcs map[string]Handler) map[string]lua.LGFunction {
	e := map[string]lua.LGFunction{}

	for name, handler := range funcs {
		e[name] = func(state *lua.LState) int {
			c := &Context{
				state: state,
			}

			err := handler(c)
			if err != nil {
				c.error(err.Error())
				return 0
			}

			return c.pushed
		}
	}

	return e
}
