package binder

import (
	"github.com/yuin/gopher-lua"
)

// Loader is basic loader object
type Loader struct {
	funcs   map[string]Handler
	modules []*Module
	tables  []*Table
}

// Func assign handler with specified alias
func (l *Loader) Func(name string, handler Handler) {
	l.funcs[name] = handler
}

// Module creates new module and returns it
func (l *Loader) Module(name string) *Module {
	m := &Module{
		name:   name,
		fields: map[string]lua.LValue{},
		funcs:  map[string]Handler{},
	}

	l.modules = append(l.modules, m)
	return m
}

// Table creates new table and returns it
func (l *Loader) Table(name string) *Table {
	t := &Table{
		name:    name,
		static:  map[string]Handler{},
		dynamic: map[string]Handler{},
	}

	l.tables = append(l.tables, t)
	return t
}

func (l *Loader) load(s *lua.LState) {
	f := exports(l.funcs)
	for name, fn := range f {
		s.SetGlobal(name, s.NewFunction(fn))
	}

	for _, m := range l.modules {
		m.state = s
		m.load()
	}

	for _, t := range l.tables {
		t.state = s
		t.load()
	}
}
