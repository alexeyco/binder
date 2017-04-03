package binder

import "github.com/yuin/gopher-lua"

type Table struct {
	name    string
	state   *lua.LState
	static  map[string]Handler
	dynamic map[string]Handler
}

func (t *Table) Static(name string, handler Handler) {
	t.static[name] = handler
}

func (t *Table) Dynamic(name string, handler Handler) {
	t.dynamic[name] = handler
}

func (t *Table) load() {
	mt := t.state.NewTypeMetatable(t.name)
	t.state.SetGlobal(t.name, mt)

	f := exports(t.static)
	for name, fn := range f {
		t.state.SetField(mt, name, t.state.NewFunction(fn))
	}

	t.state.SetField(mt, "__index", t.state.SetFuncs(t.state.NewTable(), exports(t.dynamic)))
}
