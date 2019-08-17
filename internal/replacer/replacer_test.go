package replacer

import (
	"reflect"
	"testing"
)

func Test_prefixReplace_Replace(t *testing.T) {
	tests := []struct {
		name      string
		oldPrefix string
		changeTo  string
		old       string
		want      Variant
	}{
		{
			name:      "mismatch",
			oldPrefix: "gen/",
			changeTo:  "gitlab.stageoffice.ru/UCS-COMMON/schema/",
			old:       "github.com/sirkon/message",
			want:      Nothing{},
		},
		{
			name:      "match",
			oldPrefix: "gen/",
			changeTo:  "gitlab.stageoffice.ru/UCS-COMMON/schema/",
			old:       "gen/marker",
			want:      Replacement("gitlab.stageoffice.ru/UCS-COMMON/schema/marker"),
		},
		{
			name:      "full-match",
			oldPrefix: "github.com/sirkon/goproxy/",
			changeTo:  "github.com/sirkon/goproxy/v2/",
			old:       "github.com/sirkon/goproxy",
			want:      Replacement("github.com/sirkon/goproxy/v2"),
		},
		{
			name:      "mismatch-2",
			oldPrefix: "gen/",
			changeTo:  "gitlab.stageoffice.ru/UCS-COMMON/schema/",
			old:       "gene/marker",
			want:      Nothing{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Prefix(tt.oldPrefix, tt.changeTo)
			if got := p.Replace(tt.old); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replace() = %#v, want %v", got, tt.want)
			}
		})
	}
}

func Test_regexpReplace_Replace(t *testing.T) {
	tests := []struct {
		name          string
		from          string
		to            string
		old           string
		wantInitError bool
		want          Variant
	}{
		{
			name:          "sample",
			from:          `^gen/(.*)$`,
			to:            "gitlab.stageoffice.ru/UCS-COMMON/schema/$1",
			old:           "gen/caddy/marker",
			wantInitError: false,
			want:          Replacement("gitlab.stageoffice.ru/UCS-COMMON/schema/caddy/marker"),
		},
		{
			name:          "ill-from",
			from:          "^gen/(.*$",
			wantInitError: true,
		},
		{
			name: "mismatch-2",
			from: `^gen/(.*)$`,
			to:   "gitlab.stageoffice.ru/UCS-COMMON/schema/$1",
			old:  "gene/marker",
			want: Nothing{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Regexp(tt.from, tt.to)
			if err != nil {
				if tt.wantInitError {
					return
				}
				t.Error(err)
				return
			}
			if tt.wantInitError {
				t.Errorf("Regexp had to return a error on invalid regexp `%s`", tt.from)
			}
			if got := r.Replace(tt.old); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replace() = %#v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replacerVersioned_Replace(t *testing.T) {
	type fields struct {
		base       string
		importHead string
		curVersion int
		newVersion int
	}
	tests := []struct {
		name          string
		base          string
		jump          int
		args          string
		wantInitError bool
		want          Variant
	}{
		{
			name: "from-suffixless",
			base: "github.com/user/project",
			jump: 2,
			args: "github.com/user/project/data",
			want: Replacement("github.com/user/project/v3/data"),
		},
		{
			name: "full-match",
			base: "github.com/user/project",
			jump: 1,
			args: "github.com/user/project",
			want: Replacement("github.com/user/project/v2"),
		},
		{
			name: "ignore-suffix-with-suffixless",
			base: "github.com/user/project",
			jump: 2,
			args: "github.com/user/project/v2/data",
			want: Nothing{},
		},
		{
			name: "jump-from-versioned",
			base: "github.com/user/project/v2",
			jump: 1,
			args: "github.com/user/project/v2/data",
			want: Replacement("github.com/user/project/v3/data"),
		},
		{
			name: "no-change",
			base: "github.com/user/project",
			jump: 2,
			args: "example.com/project",
			want: Nothing{},
		},
		{
			name:          "init-error",
			base:          "github.com/user/project/v1",
			jump:          2,
			args:          "",
			wantInitError: true,
			want:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Versioned(tt.base, tt.jump)
			if err != nil {
				if tt.wantInitError {
					return
				}
				t.Error(err)
				return
			}
			if got := r.Replace(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}
