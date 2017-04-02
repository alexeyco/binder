package main

import (
	"log"

	"github.com/alexeyco/binder"
)

func main() {
	// Based on https://github.com/yuin/gopher-lua#creating-a-module-by-go

	b := binder.New()
	math := b.Module("mymodule")
	math.Func("double", func(c *binder.Context) error {
		x := c.Param(1).Int()
		c.Push().Any(x*2)
		return nil
	})

	if err := b.DoFile("./example.lua"); err != nil {
		log.Fatalln(err)
	}
}
