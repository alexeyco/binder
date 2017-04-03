# Binder

[![Travis](https://img.shields.io/travis/alexeyco/binder.svg)](https://travis-ci.org/alexeyco/binder)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/alexeyco/binder/badge.svg?branch=master)](https://coveralls.io/github/alexeyco/binder?branch=master)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/alexeyco/binder)](https://goreportcard.com/report/github.com/alexeyco/binder)&nbsp;[![GoDoc](https://godoc.org/github.com/alexeyco/binder?status.svg)](https://godoc.org/github.com/alexeyco/binder)

Package binder allows to easily bind to Lua. Based on [gopher-lua](https://github.com/yuin/gopher-lua).

Write less, do more.

## Examples

### Functions
```go
package main

import (
	"errors"
	"log"

	"github.com/alexeyco/binder"
)

func main() {
	b := binder.New(binder.Options{
		SkipOpenLibs: true,
	})

	b.Func("log", func(c *binder.Context) error {
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
```

### Modules
```go
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

	if err := b.ExecString(`
		local r = require('reverse')

		print(r.string('ABCDEFGHIJKLMNOPQRSTUFVWXYZ'))
	`); err != nil {
		log.Fatalln(err)
	}
}

```
### Tables
```go
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

	if err := b.ExecString(`
		local p = person.new('Steeve')
		print(p:name())

		p:name('Alice')
		print(p:name())
	`); err != nil {
		log.Fatalln(err)
	}
}
```

### Configuration
Soon...

## License
```
MIT License

Copyright (c) 2017 Alexey Popov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
