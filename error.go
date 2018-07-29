package binder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/fatih/color"
	"github.com/yuin/gopher-lua"
)

const (
	// errorLinesBefore lines before problem line
	errorLinesBefore = 5

	// errorLinesBefore lines after problem line
	errorLinesAfter = 5
)

type errSourceHandler func(problem int) *source

type line struct {
	number  int
	problem bool
	code    string
}

const (
	errBgColor = color.BgRed
	errFgColor = color.FgWhite

	numberFgColor = color.FgHiBlack
	numberBgColor = color.BgBlack
)

func (l *line) str(numbers, len int) string {
	numberBg := numberBgColor

	line := l.codeString(len)
	if l.problem {
		numberBg = errBgColor
		line = color.New(errBgColor, errFgColor).Sprintf("%s", line)
	} else {
		line = highlight(line)
	}

	number := color.New(numberBg, numberFgColor).Sprintf("%s ", l.numberString(numbers))

	return number + line
}

func (l *line) numberString(pad int) string {
	output := fmt.Sprintf("%d", l.number)
	pad -= len(output)

	if pad <= 0 {
		return output
	}

	return fmt.Sprintf("%s%s", strings.Repeat(" ", pad), output)
}

func (l *line) codeString(pad int) string {
	output := l.code
	pad -= len(output)

	if pad <= 0 {
		return output
	}

	return fmt.Sprintf("%s%s", output, strings.Repeat(" ", pad))
}

type source struct {
	lines []*line
}

func (s *source) problem() string {
	output := make([]string, len(s.lines))
	numbers, len := s.maxLengths()

	for i, l := range s.lines {
		output[i] = l.str(numbers, len)
	}

	return strings.Join(output, "\n")
}

func (s *source) maxLengths() (int, int) {
	n := 0
	l := 0

	for _, line := range s.lines {
		nl := len(fmt.Sprintf("%d", line.number))
		if nl > n {
			n = nl
		}

		ll := len(line.code)
		if ll > l {
			l = ll
		}
	}

	return n, l
}

func newSource(code string, problem int) *source {
	l := strings.Split(code, "\n")
	count := len(l)
	lines := make([]*line, count)

	for i, v := range l {
		lines[i] = &line{
			number:  i + 1,
			problem: i+1 == problem,
			code:    v,
		}
	}

	length := errorLinesBefore + errorLinesAfter + 1

	if length <= count {
		before := problem - errorLinesBefore
		after := problem + errorLinesAfter

		if before < 0 {
			before = 0
			after = before + length
		}

		if after > count-1 {
			after = count - 1
			before = after - length
		}

		lines = lines[before-1 : after]
	}

	return &source{
		lines: lines,
	}
}

// Error error object
type Error struct {
	error   string
	problem int
	source  *source
}

// Error returns error string
func (e *Error) Error() string {
	if e.problem >= 0 {
		return fmt.Sprintf("Line %d: %s", e.problem, e.error)
	}

	return e.error
}

// Source returns problem source code as string
func (e *Error) Source() string {
	return e.source.problem()
}

// Print prints problem source code
func (e *Error) Print() {
	color.New(color.FgHiRed).Println(e.Error())
	fmt.Println()
	fmt.Println(e.Source())
}

func newError(err error, h errSourceHandler) error {
	switch err.(type) {
	case *lua.ApiError:
		e := err.(*lua.ApiError)
		p := strings.Split(e.Object.String(), ":")
		c := len(p)

		var (
			problem = -1
			text    = e.Error()
		)

		if c > 1 {
			if v, err := strconv.Atoi(p[1]); err == nil {
				problem = v
			}
		}

		if c > 2 {
			text = strings.Trim(p[2], " ")
		}

		return &Error{
			error:   text,
			problem: problem,
			source:  h(problem),
		}
	}

	return err
}

func highlight(s string) string {
	var buf bytes.Buffer
	if err := quick.Highlight(&buf, s, "lua", "terminal", "solarized-dark"); err != nil {
		return s
	}

	b, err := ioutil.ReadAll(&buf)
	if err != nil {
		return s
	}

	return string(b)
}
