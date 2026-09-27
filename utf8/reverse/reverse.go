//go:build !solution

package reverse
import "strings"

func byteCount(b byte) int {
	// 0xxxxxxx
	// 110xxxxx
	// 1110xxxx
	// 11110xxx
	b = b >> 4
	if b >= 0 && b < 8 {
		return 1
	}
	if b == 12 || b == 13 {
		return 2
	}
	if b == 14 {
		return 3
	}
	if b == 15 {
		return 4
	}
	return 0
}



func Reverse(input string) string {
	builder := &strings.Builder{}

	for i := len(input) - 1; i >= 0; i-- {
        //find first non continuation byte(10xxxxxx) index and set `i` to that index
		for input[i] >> 6 == 2 && i > 0{
			i -= 1
		}

		// get byte sequence length `n` (1-4)
		// `n` < 1 means invalid sequence
		// `i` + n 
        n := byteCount(input[i])
		valid := n > 0 && i + n <= len(input)                                                                                     
        for k := 1; k < n && valid; k++ {                                                                                       
          if input[i + k] >> 6 != 2 {                                                                                               
            valid = false                                                                                                       
          }
        }
		
		if !valid {
			builder.WriteRune('\uFFFD')
		} else {
			for k := i; k < i + n; k++ {
				builder.WriteByte(input[k])
			}
		}
	}
	return builder.String()
}

