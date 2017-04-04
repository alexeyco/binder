package binder

import (
	"github.com/yuin/gopher-lua"
)

// Handler is binder function handler
type Handler func(*Context) error

// Binder is a binder... that's all
type Binder struct {
	*Loader
	state   *lua.LState
	loaders []*Loader
}

// Load apply Loader
func (b *Binder) Load(loader *Loader) {
	b.loaders = append(b.loaders, loader)
}

// DoString runs lua script string
func (b *Binder) DoString(s string) error {
	b.load()
	return b.state.DoString(s)
}

// DoFile runs lua script file
func (b *Binder) DoFile(f string) error {
	b.load()
	return b.state.DoFile(f)
}

func (b *Binder) load() {
	loaders := append([]*Loader{b.Loader}, b.loaders...)

	for _, l := range loaders {
		l.load(b.state)
	}
}

// New returns new binder instance
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

	b := &Binder{
		state:   s,
		loaders: []*Loader{},
	}
	b.Loader = NewLoader()

	return b
}

// NewLoader returns new loader
func NewLoader() *Loader {
	return &Loader{
		funcs:   map[string]Handler{},
		modules: []*Module{},
		tables:  []*Table{},
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
