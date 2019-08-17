package replacer

import (
	"fmt"
)

type Variant interface {
	variant()
}

var _ Variant = Nothing{}

type Nothing struct{}

func (Nothing) variant() {}

var _ Variant = Replacement("")
var _ fmt.Stringer = Replacement("")

type Replacement string

func (r Replacement) String() string { return string(r) }

func (r Replacement) variant() {}
