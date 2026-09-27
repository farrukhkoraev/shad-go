//go:build !solution

package speller
import "fmt"

func digitsOf(n int64) []int {
	digits := []int{}
	
	for n > 0 {
		digits = append(digits, n % 10)
		n /= 10
	}
	
	return digits
}
func Spell(n int64) string {
	return ""
}
