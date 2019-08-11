package palindrome

import (
	"math/rand"
	"unicode"
)

func IsPalindrome(s string) bool {

	var letters []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	for i := range letters {
		if len(letters)/2+1 > i {
			if letters[i] != letters[len(letters)-1-i] {
				return false
			}
		}
	}
	return true
}

func randomizationPalindrome(rnd *rand.Rand) string {
	n := rnd.Intn(35)
	runes := make([]rune, n)

	for i := 0; i < (n+1)/2; i++ {
		r := rune(rnd.Intn(0x1000))
		runes[i] = r
		runes[n-1-i] = r
	}
	return string(runes)
}

func randomizationNonPalindrome(rnd *rand.Rand) string {
	n := rnd.Intn(35)
	runes := make([]rune, n)

	for i := 0; i < (n+1)/2; i++ {
		r := rune(rnd.Intn(0x1000))
		runes[i] = r
		r = rune(rnd.Intn(0x1000))
		runes[n-i-1] = r
	}
	return string(runes)
}
