//go:build !solution

package spacecollapse
import "strings"


func isSpaceRune(r rune) bool{
	return (r == ' ' || r == '\t' || r == '\n' || r == '\r') 
}

func CollapseSpaces(input string) string {
	builder := &strings.Builder {}
	
	wasSpace := false
	for _, k := range input {
		if isSpaceRune(k) {
			if wasSpace {
				continue
			} else {
				wasSpace = true
				builder.WriteRune(' ')
			}
		} else {
			wasSpace = false
			builder.WriteRune(k)
		}
	} 
	
	return builder.String()
}
