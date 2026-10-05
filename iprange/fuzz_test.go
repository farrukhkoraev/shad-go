package iprange

import "testing"

func FuzzParseList(f *testing.F) {
	testCases := []string {
		"",
		"192.168.1.1",
		"10.0.0.0/24",
		"192.168.1.*",
		"192.168.1.10-20",
		"192.168.1.1, 10.0.0.0/8",
	}
	for _, c := range testCases {
		f.Add(c)
	}

	f.Fuzz(func(t *testing.T, s string) {
		ParseList(s)
	})
}
