package main

import (
	"log"
	"errors"

	"github.com/alexeyco/binder"
)

type Person struct {
	n string
}

func (p *Person) Name() string {
	return p.n
}

func (p *Person) SetName(name string) {
	p.n = name
}

func main() {
	// Based on https://github.com/yuin/gopher-lua#user-defined-types

	b := binder.New()

	t := b.Table("person")
	t.Static("new", func(c *binder.Context) error {
		p := &Person{}
		p.SetName(c.Param(1).String())
		c.Push().Data(p)

		return nil
	})

	t.Method("name", func(c *binder.Context) error {
		if person, ok := c.Data(1).Value().(*Person); ok {
			if c.Top() == 1 {
				c.Push().Any(person.Name())
			} else {
				person.SetName(c.Param(2).String())
			}

			return nil
		}

		return errors.New("person expected")
	})

	if err := b.DoFile("./example.lua"); err != nil {
		log.Fatalln(err)
	}
}
