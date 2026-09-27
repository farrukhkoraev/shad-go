//go:build !solution

package externalsort

import (
	"io"
	"os"
	"sort"
	"container/heap"

)

type Item struct {
	line string
	reader LineReader
}

type MinHeap []Item

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i].line < h[j].line
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(Item))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0: n-1]
	
	return x
}

type LineR struct {
	r io.Reader
	buf []byte
}

type LineW struct {
	w io.Writer
}


func (lr *LineR) ReadLine() (string, error) {
	getLine := func () (string, bool){
		for i, b := range lr.buf {
			if b == '\n' {
				line := string(lr.buf[:i])
		    	lr.buf = lr.buf[i+1:]
		    	return line, true
			}
		}
		return "", false
	}


	buf := make([]byte, 1024)
	for { 
		// is there a line in remaining bytes
		if line, ok := getLine(); ok {
			return line, nil
		}
		
    	// no line found read more
      	n, err := lr.r.Read(buf)
    	lr.buf = append(lr.buf, buf[:n]...)
    	
    	if err == io.EOF {
    		// if there is a line return otherwise return buf contents
    		if line, ok := getLine(); ok {
    			return line, nil
    		} else {
    			line := string(lr.buf)
    			lr.buf = nil
    			return line, io.EOF
    		}
    	} else if err != nil {
    		return "", err
    	}
	}

}

func (lw *LineW) Write(l string) error {
	_, err := lw.w.Write([]byte(l + "\n"))
	return err
}

func NewReader(r io.Reader) LineReader {
	return &LineR{r: r}
}

func NewWriter(w io.Writer) LineWriter {
	return &LineW{w: w}
}

func Merge(w LineWriter, readers ...LineReader) error {

	// initialize heap
	h := &MinHeap{}
	for _, r := range readers {
		line, err := r.ReadLine()
		if err == io.EOF {
			if line != "" {
				*h = append(*h, Item{line, r})
			}
			continue
		} else if err != nil {
			return err
		}
		*h = append(*h, Item{line, r})
	}
	heap.Init(h)

	for h.Len() > 0 {
		item := heap.Pop(h).(Item)
		if err := w.Write(item.line); err != nil {
			return err
		}

		line, err := item.reader.ReadLine()
		if err == io.EOF {
			if line != "" {
				heap.Push(h, Item{line, item.reader})
			}
		} else if err != nil {
			return err
		} else {
			heap.Push(h, Item{line, item.reader})
		}
	}
	return nil
}

func sortFile(name string) error {
	lines := []string{}
    f, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	lr, lw  := NewReader(f), NewWriter(f)

	// read all the lines
	for {
		l, err := lr.ReadLine()
 
		if err == io.EOF {
			if l != "" {
				lines = append(lines, l)
			}
			break
		} else if err != nil {
			return err
		}
		lines = append(lines, l)
	}
	
	// write back sorted to the file
	sort.Strings(lines)
	f.Seek(0, io.SeekStart)
	f.Truncate(0)
	for _, line := range lines {
		err = lw.Write(line)
		if err != nil {
			return err
		}
	}
    return f.Close()
}


func Sort(w io.Writer, in ...string) error {
	// sort each files seperately
	for _, name := range in {
		err := sortFile(name)
		if err != nil {
			return err
		}
		
	}
	
	files := []*os.File{}
	readers := []LineReader{}
	
	for _, name := range in {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		
		files = append(files, f)
		readers = append(readers, NewReader(f))
	}

	writer := NewWriter(w)
	err := Merge(writer, readers...)

	for _, f := range files {
		f.Close()
	}
	
	return err
}
