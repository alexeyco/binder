package binder

import "testing"

func TestTable_Static(t *testing.T) {
	b := New()
	tbl := b.Table("test")
	tbl.Static("name", func(c *Context) error {
		return nil
	})

	if len(tbl.static) != 1 {
		t.Error("Wrong static length in table")
	}
}

func TestTable_Dynamic(t *testing.T) {
	b := New()
	tbl := b.Table("test")
	tbl.Dynamic("name", func(c *Context) error {
		return nil
	})

	if len(tbl.dynamic) != 1 {
		t.Error("Wrong dynamic length in table")
	}
}
