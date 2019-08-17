package replacer

import (
	"fmt"
	"path"
	"strings"
)

// Versioned a replaced that cares about major version suffixes
func Versioned(base string, jump int) (Replacer, error) {
	base = strings.TrimRight(base, "/")
	suffixless, end := path.Split(base)

	var sf Suffix
	var curVersion int
	if ok, _ := sf.Extract(end); ok {
		if sf.Major < 2 {
			return nil, fmt.Errorf("major version suffixes for versions lesser than 2 is a bad tone")
		}
		curVersion = sf.Major
	} else {
		suffixless = base
	}

	if curVersion == 0 {
		jump++
	}

	return &replacerVersioned{
		base:       strings.TrimRight(base, "/") + "/",
		importHead: path.Join(strings.TrimRight(suffixless, "/"), fmt.Sprintf("v%d", curVersion+jump)),
		curVersion: curVersion,
		newVersion: curVersion + jump,
	}, nil
}

type replacerVersioned struct {
	base       string
	importHead string
	curVersion int
	newVersion int
}

// Replace - take jump = 1 for instance. In this case this will replace
//   github.com/user/project    => github.com/user/project/v2 with base github.com/user/project
//   github.com/user/project/v3 =>  github.com/user/project/v4 with base github.com/user/project/v3
//   github.com/user/project/v2 => github.com/user/project/v2 with base github.com/user/project
func (r *replacerVersioned) Replace(old string) Variant {
	if !strings.HasPrefix(old, r.base) {
		if old == strings.TrimRight(r.base, "/") {
			return Replacement(r.importHead)
		}
		return Nothing{}
	}

	rest := old[len(r.base):]

	// it can be a greater major version import in case of suffixless base
	if r.curVersion == 0 {
		var sf Suffix
		if ok, _ := sf.Extract(rest); ok {
			return Nothing{}
		}
	}

	return Replacement(path.Join(r.importHead, rest))
}
