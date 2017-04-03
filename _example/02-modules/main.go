package main

import (
	"errors"
	"log"

	"github.com/alexeyco/binder"
)

func main() {
	b := binder.New()

	m := b.Module("reverse")
	m.Func("string", func(c *binder.Context) error {
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

		print(r.string('ABCDEFGHIJKLMNOPQRSTUFVWXYZ'))
	`); err != nil {
		log.Fatalln(err)
	}
}
