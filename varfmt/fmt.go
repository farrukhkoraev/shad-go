//go:build !solution

package varfmt
import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	runes := []rune(format)
	implicitIndex := 0
	builder := &strings.Builder{}

	for i:=0; i < len(runes); i++ {
		if runes[i] == '{' {
			j := i
			for runes[j] != '}' {
				j += 1
			}
			// {}
			 if j - i == 1 {
				builder.WriteString(fmt.Sprintf("%v", args[implicitIndex]))
				implicitIndex += 1
			 } else { //{1} ... {n}
				n, _ := strconv.Atoi(string(runes[i+1:j]))
			    builder.WriteString(fmt.Sprintf("%v", args[n]))
	     		implicitIndex += 1
			}
			
			i = j
			
		} else {
			builder.WriteRune(runes[i])
		}
	}
	return builder.String()
}
