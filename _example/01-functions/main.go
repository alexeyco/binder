package main

import (
	"log"

	"github.com/alexeyco/binder"
)

func main() {
	// Based on https://github.com/yuin/gopher-lua#user-defined-types

	b := binder.New()

	b.Func("multiply", func(c *binder.Context) error {
		x := c.Param(1).Int()
		y := c.Param(2).Int()

		c.Push().Any(x*y)

		return nil
	})

	if err := b.DoFile("./example.lua"); err != nil {
		log.Fatalln(err)
	}
}
