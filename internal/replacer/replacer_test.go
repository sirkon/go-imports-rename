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
			name:      "mismatch-2",
			oldPrefix: "gen/",
			changeTo:  "gitlab.stageoffice.ru/UCS-COMMON/schema/",
			old:       "gene/marker",
			want:      Nothing{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &prefixReplace{
				oldPrefix: tt.oldPrefix,
				changeTo:  tt.changeTo,
			}
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
			r, err := RegexpReplace(tt.from, tt.to)
			if err != nil {
				if tt.wantInitError {
					return
				}
				t.Error(err)
				return
			}
			if tt.wantInitError {
				t.Errorf("RegexpReplace had to return a error on invalid regexp `%s`", tt.from)
			}
			if got := r.Replace(tt.old); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replace() = %#v, want %v", got, tt.want)
			}
		})
	}
}
