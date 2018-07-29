package binder

import (
	"fmt"
	"strconv"
	"strings"

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

func (l *line) String(numbers int) string {
	numberBg := numberBgColor

	line := l.code
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

type source struct {
	lines []*line
}

func (s *source) problem() string {
	output := make([]string, len(s.lines))
	numbers := s.numbersMaxLength()

	for i, l := range s.lines {
		output[i] = l.String(numbers)
	}

	return strings.Join(output, "\n")
}

func (s *source) numbersMaxLength() int {
	l := 0
	for _, line := range s.lines {
		length := len(fmt.Sprintf("%d", line.number))
		if length > l {
			l = length
		}
	}

	return l
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

		lines = lines[before:after]
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

		break
	}

	return err
}

func highlight(s string) string {
	return s
}
