package palindrome

import (
	"math/rand"
	"testing"
	"time"
)

func TestIsPalindrome(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test 1 is palindrome", args{"aktka"}, true},
		{"test 2 is palindrome", args{"youyou"}, true},
		{"test 3 is palindrome", args{"111111"}, true},
		{"test 4 is palindrome", args{"Noelle Eve Elleon"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.args.s); got != tt.want {
				t.Errorf("IsPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrenchPalindrome(t *testing.T) {
	if !IsPalindrome("été") {
		t.Error(`IsPalindrome("été") = false`)
	}
}
func TestCanalPalindrome(t *testing.T) {
	input := "A man, a plan, a canal: Panama"
	if !IsPalindrome(input) {
		t.Errorf(`IsPalindrome(%q) = false`, input)
	}
}

func TestIsPalindrome2(t *testing.T) {
	var inputs = []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"kayak", true},
		{"deserts", true},
	}
	for _, input := range inputs {
		if got := IsPalindrome(input.input); got != input.want {
			t.Errorf("IsPalindrom(%q) = %v", input.input, got)
		}
	}
}

func Test_randomization(t *testing.T) {
	type args struct {
		rnd *rand.Rand
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test 1 random", args{rand.New(rand.NewSource(time.Now().UTC().UnixNano()))}, true},
		{"test 2 random", args{rand.New(rand.NewSource(time.Now().UTC().UnixNano()))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := randomizationPalindrome(tt.args.rnd)
			if got := IsPalindrome(input); got != tt.want {
				t.Errorf("IsPalindrom(%q) = %v, want %v", input, got, tt.want)
			}
		})
	}
}

func Test_randomizationNonPalindrome(t *testing.T) {
	type args struct {
		rnd *rand.Rand
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test 1 random", args{rand.New(rand.NewSource(time.Now().UTC().UnixNano()))}, true},
		{"test 2 random", args{rand.New(rand.NewSource(time.Now().UTC().UnixNano()))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := randomizationNonPalindrome(tt.args.rnd)
			if got := IsPalindrome(input); got != tt.want {
				t.Errorf("IsPalindrom(%q) = %v, want %v", input, got, tt.want)
			}
		})
	}
}
