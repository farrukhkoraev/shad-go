//go:build !solution

package fileleak

import (
	"os"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func openFiles() (map[string]string, error) {
	result := map[string]string{}
	
	files, err := os.ReadDir("/proc/self/fd/")
	if err != nil {
		return result, nil
	}
	for _, f := range files {
		v, err := os.Readlink("/proc/self/fd/" + f.Name())
		if err != nil {
			continue
		}
		result[f.Name()] = v
	}	
	return result, nil
}

func VerifyNone(t testingT) {
	before, _ := openFiles()
	cleanup := func() {
		after, _ := openFiles()
		count := 0
		for k, v  := range after {
			if v != before[k] {
				count += 1
			}
		}
		if count > 0 {
			t.Errorf("Leak found: %d files open", count)
		}
	}
	t.Cleanup(cleanup)
}
