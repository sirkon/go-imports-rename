package replacer

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Replacer interface {
	Replace(old string) Variant
}

type prefixReplace struct {
	oldPrefix	string
	changeTo	string
}

func PrefixReplace(old, changeTo string) Replacer {
	return &prefixReplace{
		oldPrefix:	old,
		changeTo:	changeTo,
	}
}

func (p *prefixReplace) Replace(old string) Variant {
	if !strings.HasPrefix(old, p.oldPrefix) {
		return Nothing{}
	}

	return Replacement(p.changeTo + old[len(p.oldPrefix):])
}

var _ Replacer = &regexpReplace{}

type regexpReplace struct {
	from	*regexp.Regexp
	to	string
}

func RegexpReplace(from string, to string) (Replacer, error) {
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
