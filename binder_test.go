package binder

import "testing"

func TestBinder_Func(t *testing.T) {
	b := New()
	b.Func("test", func(c *Context) error {
		return nil
	})

	if len(b.funcs) != 1 {
		t.Error("Wrong binder funcs slice length")
	}
}

func TestBinder_Table(t *testing.T) {
	b := New()
	b.Table("test")

	if len(b.tables) != 1 {
		t.Error("Wrong binder tables slice length")
	}
}

func TestBinder_Module(t *testing.T) {
	b := New()
	b.Module("test")

	if len(b.modules) != 1 {
		t.Error("Wrong binder modules slice length")
	}
}

func TestNew(t *testing.T) {
	b := New()
	if b == nil || b.state == nil || len(b.funcs) != 0 || len(b.modules) != 0 || len(b.tables) != 0 {
		t.Error("Wrong binder instance")
	}
}
