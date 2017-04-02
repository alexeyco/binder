package main

import (
	"log"

	"github.com/alexeyco/binder"
)

func main() {
	// Based on https://github.com/yuin/gopher-lua#calling-go-from-lua

	b := binder.New()

	b.Func("double", func(c *binder.Context) error {
		x := c.Param(1).Int()

		c.Push().Value(x*2)

		return nil
	})

	if err := b.DoFile("./example.lua"); err != nil {
		log.Fatalln(err)
	}
}
