package main

import (
	"errors"
	"log"

	"github.com/alexeyco/binder"
)

type Person struct {
	Name string
}

func main() {
	b := binder.New()

	t := b.Table("person")
	t.Static("new", func(c *binder.Context) error {
		if c.Top() == 0 {
			return errors.New("need arguments")
		}
		n := c.Arg(1).String()

		c.Push().Data(&Person{n}, "person")
		return nil
	})

	t.Dynamic("name", func(c *binder.Context) error {
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
		print(p:name())

		p:name('Alice')
		print(p:name())
	`); err != nil {
		log.Fatalln(err)
	}
}
