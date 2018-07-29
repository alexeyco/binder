package binder

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorSourceSmall_Func(t *testing.T) {
	b := getBinder()

	if err := b.DoString(`
local p = person.new('Steeve')
print(p:email())
    `); err != nil {
		switch err.(type) {
		case *Error:
			e := err.(*Error)
			s := strings.Split(e.Source(), "\n")
			l := len(s)

			if l != 4 {
				t.Errorf("Source must have %d lines, received %d", 4, l)
			}

			break
		default:
			t.Error("Must return error", err)
		}
	}
}

func getBinder() *Binder {
	b := New()
	tbl := b.Table("person")
	tbl.Static("new", func(c *Context) error {
		if c.Top() == 0 {
			return errors.New("need arguments")
		}
		n := c.Arg(1).String()

		c.Push().Data(&Person{n}, "person")
		return nil
	})

	tbl.Dynamic("name", func(c *Context) error {
		p, ok := c.Arg(1).Data().(*Person)
		if !ok {
			return errors.New("person expected")
		}

		if c.Top() == 1 {
			c.Push().String(p.Name)
		} else {
			p.Name = c.Arg(2).String()
		}

		return nil
	})

	return b
}
