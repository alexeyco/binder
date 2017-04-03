package binder

import (
	"github.com/yuin/gopher-lua"
)

type Handler func(*Context) error

type Binder struct {
	state   *lua.LState
	funcs   map[string]Handler
	modules []*Module
	tables  []*Table
}

func (b *Binder) EnableBuiltins() {

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

func (b *Binder) Table(name string) *Table {
	t := &Table{
		name:    name,
		state:   b.state,
		static:  map[string]Handler{},
		dynamic: map[string]Handler{},
	}

	b.tables = append(b.tables, t)
	return t
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

	for _, t := range b.tables {
		t.load()
	}
}

func New(opts ...Options) *Binder {
	options := []lua.Options{}

	if len(opts) > 0 {
		o := opts[0]

		options = append(options, lua.Options{
			CallStackSize:       o.CallStackSize,
			RegistrySize:        o.RegistrySize,
			SkipOpenLibs:        o.SkipOpenLibs,
			IncludeGoStackTrace: o.IncludeGoStackTrace,
		})
	}

	s := lua.NewState(options...)
	defer s.Close()

	if len(opts) > 0 && opts[0].SkipOpenLibs {
		type lib struct {
			Name string
			Func lua.LGFunction
		}

		libs := []lib{
			{lua.LoadLibName, lua.OpenPackage},
			{lua.BaseLibName, lua.OpenBase},
			{lua.TabLibName, lua.OpenTable},
		}

		for _, l := range libs {
			s.Push(s.NewFunction(l.Func))
			s.Push(lua.LString(l.Name))
			s.Call(1, 0)
		}
	}

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
