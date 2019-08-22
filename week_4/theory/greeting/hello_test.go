package main

import (
	"fmt"
	"testing"
)

func Test_hello(t *testing.T) {
	type args struct {
		user string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "empty name", args: args{user: ""}, want: "Hello Dude!"},
		{name: "not empty name ", args: args{"LuanNH"}, want: "Hello LuanNH"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := hello(tt.args.user); got != tt.want {
				t.Errorf("hello() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_hello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hello("aaaa")
	}
}

func ExampleHello() {
	fmt.Printf("%s", hello("aaa"))
	// Output:
	// Hello aa
}
