package binder

import (
	"errors"
	"testing"
)

func TestLua_Func(t *testing.T) {
	b := New(Options{
		SkipOpenLibs: true,
	})

	b.Func("sum", func(c *Context) error {
		t := c.Top()
		if t < 2 {
			return errors.New("need at least 2 arguments")
		}

		var sum float64
		for i := 1; i <= t; i++ {
			sum += c.Arg(i).Number()
		}

		c.Push().Number(sum)
		return nil
	})

	if err := b.DoString(`
		assert(sum(1, 2) == 3, '1 + 2 != 3')
		assert(sum(5, 7) == 12, '5 + 7 != 12')
		assert(sum(100, 200) == 300, '100 + 200 != 300')
	`); err != nil {
		t.Error("Error execute function", err)
	}
}

func TestLua_Module(t *testing.T) {
	b := New()

	m := b.Module("reverse")
	m.Func("string", func(c *Context) error {
		if c.Top() == 0 {
			return errors.New("need arguments")
		}

		s := c.Arg(1).String()

		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}

		c.Push().String(string(runes))
		return nil
	})

	if err := b.DoString(`
		local r = require('reverse')

		assert(r.string('ABCDE') == 'EDCBA', 'ABCDE != EDCBA')
		assert(r.string('01234') == '43210', '01234 != 43210')
	`); err != nil {
		t.Error("Error execute module", err)
	}
}

type Person struct {
	Name string
}

func TestLua_Table(t *testing.T) {
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

	if err := b.DoString(`
		local p = person.new('Steeve')

		assert(p:name() == 'Steeve', 'Steve is not Steve')

		p:name('Alice')
		assert(p:name() == 'Alice', 'Alice is not Alice')
	`); err != nil {
		t.Error("Error execute module", err)
	}
}
