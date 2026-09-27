//go:build !solution

package otp

import (
	"io"
)


type Cipher struct {
	r io.Reader
	w io.Writer
	prng io.Reader
}

func (c Cipher) Read(buf []byte) (int, error){

	n, err := c.r.Read(buf)

	if n > 0 {
		rngBuf := make([]byte, n)
	    _, _ = c.prng.Read(rngBuf)

	    for i := range n{
			buf[i] = buf[i] ^ rngBuf[i] 
		}
	}
	return n, err

}

func (c Cipher) Write(buf []byte) (int, error){

	newBuf := make([]byte, len(buf))
	c.prng.Read(newBuf)
	
	for i := range newBuf {
		newBuf[i] = newBuf[i] ^ buf[i]
	}

	return  c.w.Write(newBuf)
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return Cipher{r, nil, prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return Cipher {nil, w, prng}
}
