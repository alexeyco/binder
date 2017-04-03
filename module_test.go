package binder

import "testing"

func TestModule_String(t *testing.T) {
	b := New()
	m := b.Module("test")
	m.String("foo", "bar")

	if len(m.fields) != 1 {
		t.Error("Wrong fields length in module")
	}
}

func TestModule_Number(t *testing.T) {
	b := New()
	m := b.Module("test")
	m.Number("foo", 1)

	if len(m.fields) != 1 {
		t.Error("Wrong fields length in module")
	}
}

func TestModule_Bool(t *testing.T) {
	b := New()
	m := b.Module("test")
	m.Bool("foo", true)

	if len(m.fields) != 1 {
		t.Error("Wrong fields length in module")
	}
}

func TestModule_Func(t *testing.T) {
	b := New()
	m := b.Module("test")
	m.Func("foo", func(c *Context) error {
		return nil
	})

	if len(m.funcs) != 1 {
		t.Error("Wrong funcs length in module")
	}
}
