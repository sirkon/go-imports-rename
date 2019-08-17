package replacer

import (
	"strings"
)

type prefixReplace struct {
	oldPrefix string
	changeTo  string
}

var _ Replacer = &prefixReplace{}

func Prefix(old, changeTo string) Replacer {
	return &prefixReplace{
		oldPrefix: old,
		changeTo:  changeTo,
	}
}

func (p *prefixReplace) Replace(old string) Variant {
	if !strings.HasPrefix(old, p.oldPrefix) {
		if old == strings.TrimRight(p.oldPrefix, "/") {
			return Replacement(strings.TrimRight(p.changeTo, "/"))
		}
		return Nothing{}
	}

	return Replacement(p.changeTo + old[len(p.oldPrefix):])
}
