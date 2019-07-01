package replacer

import (
	"fmt"
)

type Variant interface {
	dontTouchMe()
}

var _ Variant = Nothing{}

type Nothing struct{}

func (Nothing) dontTouchMe()	{}

var _ Variant = Replacement("")
var _ fmt.Stringer = Replacement("")

type Replacement string

func (r Replacement) String() string	{ return string(r) }

func (r Replacement) dontTouchMe()	{}
