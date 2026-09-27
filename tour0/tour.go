//go:build !solution

package main

import (
	"fmt"
	"math"
	"strings"
	"io"
	"os"
)


func Tour() string {
	return ""
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error () string {
	return fmt.Sprintf("cannot Sqrt negative number: %.2f", e)
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	var z0, z1 = 0.0, 1.0
	
	for math.Abs(z1-z0) > 0.01 {
		z0 = z1
		z1 = z1 - (z1 * z1 - x) / (2 * z1)
	}
	return z1, nil
}

func WordCount(s string) map[string]int {
	var m = make(map[string]int)
	for _, w := range(strings.Fields(s)) {
		if c, ok := m[w]; ok {
			m[w] = c + 1
		} else {
			m[w] = 1
		}
		
	}  
	return m
}

func fibonacci () (func() int) {
	var n_1, n = 0, 1
	return  func() int {
		n_1, n = n, n_1 + n
		return n
	}
}

type MyFloat = float64

type Vertex struct {
	X, Y float64
}

type Abser interface {
	Abs() float64
} 

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X * v.X + v.Y * v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f	
}

func ScaleP(v *Vertex, f float64) {
	v.X = v.X * f
    v.Y = v.Y * f
} 

func describe(i Abser) {
	fmt.Printf("(%v, %T)\n", i, i)
}

type IPAddr [4]byte

func (ip IPAddr) String() string {
	var out = make([]string, 4)
	for i, b := range(ip) {
		out[i] = fmt.Sprintf("%d", b) 
	}
	return strings.Join(out, ".")
	
}

type MyReader struct {}

func (r MyReader) Read(b []byte) (int, error) {
	for i := range(b) {
		b[i] = 'A'
	}
	return len(b), nil
}

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(b []byte) (int, error) {
	rot13 := func (b byte) byte {
		switch {
		case b >= 'A' && b <= 'Z':
			return (b-'A'+13) % 26 + 'A'
		case b >= 'a' && b <= 'z':
			return (b-'a'+13) % 26 + 'a'
		default:
			return b
		}
	}
	n, err := r.r.Read(b)
	if err != nil {
		return n, err
	}
	
	for i := range(n) {
		b[i] = rot13(b[i])
	}
	return n, nil
} 

func main() {


	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)

	// hosts := map[string]IPAddr{
	// 	"loopback":  {127, 0, 0, 1},
	// 	"googleDNS": {8, 8, 8, 8},
	// }
	// for name, ip := range hosts {
	// 	
	// 	fmt.Printf("%v: %v\n", name, ip)
	// }
}
