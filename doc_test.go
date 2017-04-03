package binder

import (
	"errors"
	"log"
)

func ExampleBinder_Func() {
	b := New(Options{
		SkipOpenLibs: true,
	})

	b.Func("log", func(c *Context) error {
		t := c.Top()
		if t == 0 {
			return errors.New("need arguments")
		}

		l := []interface{}{}

		for i := 1; i <= t; i++ {
			l = append(l, c.Arg(i).Any())
		}

		log.Println(l...)
		return nil
	})

	if err := b.ExecString(`
		log('This', 'is', 'Lua')
	`); err != nil {
		log.Fatalln(err)
	}
}

func ExampleBinder_Module() {
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

	if err := b.ExecString(`
		local r = require('reverse')
		print(r.string('ABCDEFGHIJKLMNOPQRSTUFVWXYZ'))
	`); err != nil {
		log.Fatalln(err)
	}
}

func ExampleBinder_Table() {
	type Person struct {
		Name string
	}

	b := New()

	t := b.Table("person")
	t.Static("new", func(c *Context) error {
		if c.Top() == 0 {
			return errors.New("need arguments")
		}
		n := c.Arg(1).String()

		c.Push().Data(&Person{n}, "person")
		return nil
	})

	t.Dynamic("name", func(c *Context) error {
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

	if err := b.ExecString(`
		local p = person.new('Steeve')
		print(p:name())

		p:name('Alice')
		print(p:name())
	`); err != nil {
		log.Fatalln(err)
	}
}
