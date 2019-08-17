package parser

import (
	"testing"
)

func TestScanner_NextString(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    string
		wantErr bool
	}{
		{
			name:    "ASCII",
			scanner: NewScanner("  abcd/efg"),
			want:    "abcd/efg",
			wantErr: false,
		},
		{
			name:    "Unicode",
			scanner: NewScanner("  абвг/деё"),
			want:    "абвг/деё",
			wantErr: false,
		},
		{
			name:    "space-escape",
			scanner: NewScanner(`abcd\ dcba`),
			want:    "abcd dcba",
			wantErr: false,
		},
		{
			name:    "invalid-input",
			scanner: NewScanner("abde\redba"),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scanner.NextString()
			if (err != nil) != tt.wantErr {
				t.Errorf("NextString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_NextOperator(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    string
		wantErr bool
	}{
		{
			name:    "operator-prefix",
			scanner: NewScanner("  =>"),
			want:    "=>",
			wantErr: false,
		},
		{
			name:    "operator-increment",
			scanner: NewScanner(" ++"),
			want:    "++",
			wantErr: false,
		},
		{
			name:    "operator-add",
			scanner: NewScanner("+="),
			want:    "+=",
			wantErr: false,
		},
		{
			name:    "operator-regexp",
			scanner: NewScanner("//"),
			want:    "//",
			wantErr: false,
		},
		{
			name:    "no-operator",
			scanner: NewScanner("abdef"),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scanner.NextOperator()
			if (err != nil) != tt.wantErr {
				t.Errorf("NextOperator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextOperator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_NextInt(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    int
		wantErr bool
	}{
		{
			name:    "ok",
			scanner: NewScanner("     1234"),
			want:    1234,
			wantErr: false,
		},
		{
			name:    "empty-input",
			scanner: NewScanner("           "),
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid-input",
			scanner: NewScanner("     1234a"),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scanner.NextInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("NextInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_NextString_ErrorOutput(t *testing.T) {
	scanner := NewScanner("                     ab\nba")
	if _, err := scanner.NextString(); err != nil {
		t.Log(scanner.FancyIndicator(0, 0))
	} else {
		t.Errorf("error expected")
	}

	scanner = NewScanner("abdef" + string([]byte{17, 29}))
	if _, err := scanner.NextString(); err != nil {
		t.Log(scanner.FancyIndicator(2, 0))
	} else {
		t.Error("error expected")
	}

	scanner = NewScanner("asdfsdfdf dsfasdfsdfsdfadsfds")
	t.Log(scanner.FancyIndicator(0, 2))
}
