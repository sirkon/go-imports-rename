package replacer

import (
	"regexp"

	"github.com/pkg/errors"
)

var _ Replacer = &regexpReplace{}

type regexpReplace struct {
	from *regexp.Regexp
	to   string
}

func Regexp(from string, to string) (Replacer, error) {
	fromRe, err := regexp.Compile(from)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid from regexp for regexp replacer")
	}
	return &regexpReplace{from: fromRe, to: to}, nil
}

func (r *regexpReplace) Replace(old string) Variant {
	if !r.from.MatchString(old) {
		return Nothing{}
	}

	return Replacement(r.from.ReplaceAllString(old, r.to))
}
