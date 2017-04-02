package binder

import (
	lua "github.com/yuin/gopher-lua"
)

type Binder struct {
	state  *lua.LState
	tables []*Table
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

func (b *Binder) DoString(s string) error {
	b.prepare()
	return b.state.DoString(s)
}

func (b *Binder) DoFile(f string) error {
	b.prepare()
	return b.state.DoFile(f)
}

func (b *Binder) prepare() {
	for _, t := range b.tables {
		t.prepare(b.state)
	}
}

func New() *Binder {
	s := lua.NewState()
	defer s.Close()

	return &Binder{
		state:  s,
		tables: []*Table{},
	}
}
