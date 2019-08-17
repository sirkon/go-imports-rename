package parser

import (
	"fmt"
	"io"
)

var _ error = ParseError{}

// ParseError error to
type ParseError struct {
	Report  string
	Details string
}

func (p ParseError) Error() string {
	return p.Report
}

// Parse parser input string into corresponding rule
func Parse(input string) (Rule, error) {
	scanner := NewScanner(input)

	scanner.trimSpaces()
	piece1, err := scanner.NextString()
	if err != nil {
		if err == io.EOF {
			return nil, ParseError{
				Report:  "missing opening import path or regexp",
				Details: scanner.FancyIndicator(1, 25),
			}
		}
		return nil, ParseError{
			Report:  err.Error(),
			Details: scanner.FancyIndicator(1, 0),
		}
	}

	preOp := scanner.Copy()
	operator, err := scanner.NextOperator()
	bolds := []interface{}{
		bold(operatorPrefix),
		bold(operatorVersionIncrement),
		bold(operatorVersionAdd),
		bold(operatorRegexp),
	}
	if err != nil {
		if err == io.EOF {
			return nil, ParseError{
				Report:  fmt.Sprintf("missing operator (one of %s, %s, %s or %s)", bolds...),
				Details: scanner.FancyIndicator(1, 2),
			}
		} else {
			return nil, ParseError{
				Report: fmt.Sprintf(
					"operator expected (one of %s, %s, %s or %s)", bolds...),
				Details: scanner.FancyIndicator(2, 0),
			}
		}
	}

	switch operator {
	case operatorPrefix:
		return processPrefix(scanner, piece1)
	case operatorVersionIncrement:
		if err := scanner.AtEnd(); err != nil {
			return nil, unwantedData(err, scanner)
		}
		return Add{
			Import: piece1,
			Jump:   1,
		}, nil
	case operatorVersionAdd:
		return processVersionAdd(scanner, piece1)
	case operatorRegexp:
		piece2, err := scanner.NextString()
		if err != nil {
			if err == io.EOF {
				return nil, ParseError{
					Report:  "missing replacement regexp",
					Details: scanner.FancyIndicator(1, 25),
				}
			}
			return nil, ParseError{
				Report:  err.Error(),
				Details: scanner.FancyIndicator(1, 0),
			}
		}
		if err := scanner.AtEnd(); err != nil {
			return nil, unwantedData(err, scanner)
		}
		return Regexp{
			From: piece1,
			To:   piece2,
		}, nil
	default:
		return nil, ParseError{
			Report:  "unsupported operator",
			Details: preOp.FancyIndicator(len(operator), 0),
		}
	}
}

func processPrefix(scanner *Scanner, from string) (Rule, error) {
	piece2, err := scanner.NextString()
	if err != nil {
		if err == io.EOF {
			return nil, ParseError{
				Report:  "missing replacement import path",
				Details: scanner.FancyIndicator(1, 25),
			}
		}
		return nil, ParseError{
			Report:  err.Error(),
			Details: scanner.FancyIndicator(1, 0),
		}
	}
	if err := scanner.AtEnd(); err != nil {
		return nil, unwantedData(err, scanner)
	}
	return Prefix{
		From: from,
		To:   piece2,
	}, nil
}

func processVersionAdd(scanner *Scanner, piece1 string) (Rule, error) {
	jump, err := scanner.NextInt()
	if err != nil {
		if err == io.EOF {
			return nil, ParseError{
				Report:  "missing version jump value",
				Details: scanner.FancyIndicator(1, 4),
			}
		}
		return nil, ParseError{
			Report:  "version jump value expected",
			Details: scanner.FancyIndicator(1, 0),
		}
	}
	if err := scanner.AtEnd(); err != nil {
		return nil, unwantedData(err, scanner)
	}
	return Add{
		Import: piece1,
		Jump:   jump,
	}, nil
}

func unwantedData(err error, scanner *Scanner) ParseError {
	return ParseError{
		Report:  err.Error(),
		Details: scanner.FancyIndicator(100000000, 0),
	}
}
