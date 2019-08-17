package parser

// Rule abstract rule description
type Rule interface {
	rule()
}

var _ Rule = Prefix{}

// Prefix rule description
type Prefix struct {
	From string
	To   string
}

func (Prefix) rule() {}

var _ Rule = Add{}

// Add (and increment) rule description
type Add struct {
	Import string
	Jump   int
}

func (Add) rule() {}

var _ Rule = Regexp{}

// Regexp rule description
type Regexp struct {
	From string
	To   string
}

func (Regexp) rule() {}
