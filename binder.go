package binder

import (
	lua "github.com/yuin/gopher-lua"
)

type Binder struct {
	state   *lua.LState
	funcs   map[string]Handler
	tables  []*Table
	modules []*Module
}

func (b *Binder) Func(name string, handler Handler) {
	b.funcs[name] = handler
}

func (b *Binder) Table(name string) *Table {
	t := &Table{
		Name:    name,
		static:  map[string]Handler{},
		methods: map[string]Handler{},
	}

	b.tables = append(b.tables, t)
	return t
}

func (b *Binder) Module(name string) *Module {
	m := &Module{
		Name:  name,
		funcs: map[string]Handler{},
	}

	b.modules = append(b.modules, m)
	return m
}

func (b *Binder) DoString(s string) error {
	b.prepare()
	return b.state.DoString(s)
}

func (b *Binder) DoFile(f string) error {
	b.prepare()
	return b.state.DoFile(f)
}

func (b *Binder) prepare() {
	for name, handler := range b.funcs {
		b.state.SetGlobal(name, b.state.NewFunction(func(s *lua.LState) int {
			return b.handle(s, handler)
		}))
	}

	for _, t := range b.tables {
		t.prepare(b.state)
	}

	for _, m := range b.modules {
		m.prepare(b.state)
	}
}

func (b *Binder) handle(s *lua.LState, handler Handler) int {
	c := &Context{
		e:     true,
		state: s,
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

func New() *Binder {
	s := lua.NewState()
	defer s.Close()

	return &Binder{
		state:  s,
		funcs:  map[string]Handler{},
		tables: []*Table{},
	}
}
