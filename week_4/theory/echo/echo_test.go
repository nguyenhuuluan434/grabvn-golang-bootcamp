package main

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_echo(t *testing.T) {
	type args struct {
		newline bool
		sep     string
		args    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test 1", args{newline: true, args: []string{}, sep: ""}, "\n"},
		{"test 2", args{newline: false, args: []string{}, sep: ""}, ""},
		{"test 3", args{newline: true, args: []string{"one", "two", "three"}, sep: "\t"}, "one\ttwo\tthree\n"},
		{"test 4", args{newline: true, args: []string{"a", "b", "c"}, sep: ","}, "a b c\n"},
		{"test 5", args{newline: true, args: []string{}, sep: "\n"}, ""},
	}
	t.Parallel()
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			desc := fmt.Sprintf("echo(%v,%q.%q)", tt.args.newline, tt.args.sep, tt.args.args)
			out = new(bytes.Buffer) // captured output
			if err := echo(tt.args.newline, tt.args.sep, tt.args.args); err != nil {
				t.Errorf("%s failed %v", desc, err)
				return
			}
			got := out.(*bytes.Buffer).String()
			if got != tt.want {
				t.Errorf("%s %s = %q, want %q", tt.name, desc, got, tt.want)
			}
		})
	}
}
