package srt2vtt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	_ "log"
	_ "os"
	"strings"
)

func SrtScanner(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if i < len(data)-1 && string(data[i:i+2]) == "\n\n" {
			return i + 2, data[:i+2], nil
		}
		if i < len(data)-3 && string(data[i:i+4]) == "\r\n\r\n" {
			return i + 4, data[:i+4], nil
		}
	}
	if atEOF && len(data) != 0 {
		return len(data), data, nil
	}
	// Golang v1.6
	// There is one final token to be delivered, which may be the empty string.
	// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
	// but does not trigger an error to be returned from Scan itself.
	//return 0, data, fmt.Errorf("bufio.ErrFinalToken")
	return 0, nil, nil
}

func ConvertTimeToWebVtt(t string) string {
	t = strings.Replace(t, ",", ".", -1)
	timing := strings.Split(t, " --> ")
	for i, v := range timing {
		if v[:3] == "00:" {
			timing[i] = v[3:]
		}
	}
	return strings.Join(timing, " --> ")
}

func SrtToWebVtt(l string) string {
	l = strings.Replace(l, "\r\n", "\n", -1)
	lines := strings.SplitN(l, "\n", 3)[1:]
	lines[0] = ConvertTimeToWebVtt(lines[0])
	return fmt.Sprintf("%s\n%s", lines[0], lines[1])
}

type Reader struct {
	s *bufio.Scanner
	b bytes.Buffer
	p int64
}

func NewReader(reader io.Reader) (*Reader, error) {
	r := new(Reader)
	r.p = 0
	r.s = bufio.NewScanner(reader)
	r.s.Split(SrtScanner)
	return r, nil
}

func (r *Reader) Read(p []byte) (n int, e error) {
	var buf bytes.Buffer
	if r.p == 0 {
		buf.WriteString("WEBVTT\n\n")
	} else {
		buf.Write(r.b.Bytes())
		r.b.Reset()
	}

	for buf.Len() < len(p) && r.s.Scan() {
		l := r.s.Text()
		l = SrtToWebVtt(l)
		buf.WriteString(l)
	}

	/*if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}*/

	n = copy(p, buf.Bytes())
	r.p = r.p + int64(n)
	r.b.Write(buf.Bytes()[n:])
	if n == 0 {
		return n, io.EOF
	}
	return n, nil
}
