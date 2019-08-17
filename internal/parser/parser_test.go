package parser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    string
		want    Rule
		wantErr bool
	}{
		{
			name: "prefix",
			args: "github.com/sirkon/ldetool/ => github.com/sirkon/ldetool/v2/",
			want: Prefix{
				From: "github.com/sirkon/ldetool/",
				To:   "github.com/sirkon/ldetool/v2/",
			},
			wantErr: false,
		},
		{
			name: "increment",
			args: "github.com/sirkon/ldetool ++",
			want: Add{
				Import: "github.com/sirkon/ldetool",
				Jump:   1,
			},
			wantErr: false,
		},
		{
			name: "add",
			args: "github.com/sirkon/ldetool += 5",
			want: Add{
				Import: "github.com/sirkon/ldetool",
				Jump:   5,
			},
			wantErr: false,
		},
		{
			name: "regexp",
			args: "github.com/sirkon/([^/]*)/(.*) // github/com/sirkon/ldetool/$2",
			want: Regexp{
				From: "github.com/sirkon/([^/]*)/(.*)",
				To:   "github/com/sirkon/ldetool/$2",
			},
			wantErr: false,
		},
		{
			name:    "invalid-empty",
			args:    "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-no-operator",
			args:    "import/path",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-invalid-operator",
			args:    "import/path --",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-prefix-no-replacement",
			args:    "import/path =>",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-prefix-unwanted-data",
			args:    "import/path => path unwanted",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-increment-unwanted-data",
			args:    "import/path ++ unwanted",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-add-no-jump",
			args:    "import/path += ",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-add-invalid-jump",
			args:    "import/path += a",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-add-jump-unwanted",
			args:    "import/path += 4 unwanted",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-regexp-no-replacement",
			args:    "import/path => ",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid-regexp-unwanted-replacement",
			args:    "import/path => path unwanted",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				if v, ok := err.(ParseError); ok {
					t.Error("\r" + v.Details)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
