package main

import (
	"encoding"
	"fmt"

	"github.com/sirkon/go-imports-rename/internal/parser"
)

var _ encoding.TextUnmarshaler = &RuleType{}

// RuleType describes parser type
type RuleType struct {
	Rule parser.Rule
}

func (r *RuleType) UnmarshalText(rawText []byte) error {
	text := string(rawText)
	rule, err := parser.Parse(text)
	if err != nil {
		if v, ok := err.(parser.ParseError); ok {
			return fmt.Errorf("invalid rule, %s: %s", v.Report, v.Details)
		}
		return fmt.Errorf("invalid rule, %s", err)
	}
	r.Rule = rule
	return nil
}
