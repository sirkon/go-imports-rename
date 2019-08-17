package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	operatorPrefix           = "=>"
	operatorVersionIncrement = "++"
	operatorVersionAdd       = "+="
	operatorRegexp           = "//"
)

// Scanner input scanner
type Scanner struct {
	orig []rune
	rest []rune
}

// NewScanner scanner constructor
func NewScanner(orig string) *Scanner {
	input := []rune(orig)
	return &Scanner{
		orig: input,
		rest: input,
	}
}

// Copy creates a copy of given scanner
func (s *Scanner) Copy() *Scanner {
	s.trimSpaces()
	return &Scanner{
		orig: s.orig,
		rest: s.rest,
	}
}

func (s *Scanner) trimSpaces() {
	pos := -1
	spacesOnly := true
	for i, r := range s.rest {
		pos = i
		if r != ' ' {
			spacesOnly = false
			break
		}
	}
	if spacesOnly {
		s.rest = s.rest[len(s.rest):]
		return
	}
	if pos >= 0 {
		s.rest = s.rest[pos:]
	}
}

// NextString retrieves next contiguous text
func (s *Scanner) NextString() (string, error) {
	s.trimSpaces()
	if len(s.rest) == 0 {
		return "", io.EOF
	}
	var buf bytes.Buffer
	if err := s.scanString(&buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var _ error = operatorExpected{}

type operatorExpected struct{}

func (operatorExpected) Error() string {
	return "operator expected"
}

// NextOperator returns operator (one of =>, ++, += and //)
func (s *Scanner) NextOperator() (string, error) {
	s.trimSpaces()

	if len(s.rest) == 0 {
		return "", io.EOF
	}

	rest := string(s.rest)
	if len(rest) < 2 {
		return "", operatorExpected{}
	}
	piece := rest[:2]
	switch piece {
	case operatorPrefix:
	case operatorVersionIncrement:
	case operatorVersionAdd:
	case operatorRegexp:
	default:
		return "", operatorExpected{}
	}
	s.rest = s.rest[2:]

	return piece, nil
}

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

// FancyIndicator цветной выхлоп для диагностики ошибок в правиле
func (s *Scanner) FancyIndicator(rng int, questions int) string {
	const radius = 7

	var buf bytes.Buffer

	orig := make([]rune, len(s.orig))
	copy(orig, s.orig)

	var isSpace bool
	for i, r := range orig {
		isSpace = r == ' ' || r == '\t'
		var p string
		var useHighlighting bool
		switch r {
		case '\t':
			p = `\t`
		case '\r':
			p = `\r`
		case '\n':
			p = `\n`
		default:
			if r < ' ' {
				p = fmt.Sprintf("\\%d", int(r))
			} else {
				useHighlighting = true
				p = string(r)
			}
		}
		if s.pos() <= i && i <= s.pos()+rng {
			if useHighlighting {
				p = fmt.Sprintf("\033[31;1m%s\033[0m", p)
			} else {
				p = fmt.Sprintf("\033[31m%s\033[0m", p)
			}
		} else {
			if useHighlighting {
				p = fmt.Sprintf("\033[1m%s\033[0m", p)
			}
		}
		buf.WriteString(p)
	}
	if questions > 0 {
		if !isSpace {
			buf.WriteString(" ")
		}
		buf.WriteString("\033[31m")
		buf.WriteString(strings.Repeat("?", questions))
		buf.WriteString("\033[0m")
	}
	return buf.String()
}

// pos returns current position in the input
func (s *Scanner) pos() int {
	s.trimSpaces()
	return len(s.orig) - len(s.rest)
}

func (s *Scanner) scanString(buf *bytes.Buffer) error {
	if len(s.rest) == 0 {
		return nil
	}
	switch s.rest[0] {
	case '\\':
		if len(s.rest) < 2 {
			return errors.New("space expected after \\")
		}
		if s.rest[1] != ' ' {
			s.rest = s.rest[1:]
			return errors.New("space expected after \\")
		}
		s.rest = s.rest[2:]
		buf.WriteRune(' ')
	case ' ':
		return nil
	case '\n':
		return errors.New("\\n characters are not allowed")
	case '\t':
		return errors.New("\\t characters are not allowed")
	case '\r':
		return errors.New("\\r characters are not allowed")
	default:
		if s.rest[0] < ' ' {
			return fmt.Errorf("non-printable characters are not allowed, got %d in the input", int(s.rest[0]))
		}
		buf.WriteRune(s.rest[0])
		s.rest = s.rest[1:]
	}
	return s.scanString(buf)
}

// NextInt returns int
func (s *Scanner) NextInt() (int, error) {
	s.trimSpaces()

	if len(s.rest) == 0 {
		return 0, io.EOF
	}
	res, err := s.scanInt(0)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *Scanner) scanInt(prev int) (int, error) {
	if len(s.rest) == 0 {
		return prev, nil
	}
	item := s.rest[0]
	if item == ' ' {
		return prev, nil
	}
	if !unicode.IsDigit(item) {
		return 0, fmt.Errorf("digit expected")
	}
	s.rest = s.rest[1:]
	return s.scanInt(prev*10 + int(item-'0'))
}

// AtEnd проверка, что правило вычитано
func (s *Scanner) AtEnd() error {
	s.trimSpaces()
	if len(s.rest) > 0 {
		return fmt.Errorf("unexpected data")
	}
	return nil
}
