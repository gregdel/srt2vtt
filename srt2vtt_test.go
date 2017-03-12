package srt2vtt

import (
	"bytes"
	"os"
	"testing"
)

func TestConvertTimeToWebVtt(t *testing.T) {
	patterns := map[string]string{
		"00:00:00,579 --> 00:00:02,145": "00:00.579 --> 00:02.145",
		"01:00:00,579 --> 01:00:02,145": "01:00:00.579 --> 01:00:02.145",
	}
	for k, v := range patterns {
		result := ConvertTimeToWebVtt(k)
		if result != v {
			t.Errorf("excepted %q get %q", v, result)
		}
	}
}

func TestSrtToWebVtt(t *testing.T) {
	patterns := map[string]string{
		"2\n00:00:02,147 --> 00:00:04,257\nOh, look at all\nthat chest hair.\n\n":           "00:02.147 --> 00:04.257\nOh, look at all\nthat chest hair.\n\n",
		"2\r\n00:00:02,147 --> 00:00:04,257\r\nOh, look at all\r\nthat chest hair.\r\n\r\n": "00:02.147 --> 00:04.257\nOh, look at all\nthat chest hair.\n\n",
		"2\n00:00:02,147 --> 00:00:04,257\n♪♪♪♪\n\n":                                        "00:02.147 --> 00:04.257\n♪♪♪♪\n\n",
	}
	for k, v := range patterns {
		result, err := SrtToWebVtt(k)
		if err != nil {
			t.Fatal(err)
		}
		if result != v {
			t.Errorf("excepted %q get %q", v, result)
		}
	}

}

func TestWriteTo(t *testing.T) {
	f, err := os.Open("test/test2.srt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, err := NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	n, err := r.WriteTo(buf)
	if err == nil || n != 8 {
		t.Fatalf("written %v, err %v, buf %s", n, err, buf.String())
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func testReader(t *testing.T, bufSize int) {
	fs, err := os.Open("test/test.srt")
	check(err)
	defer fs.Close()
	fv, err := os.Open("test/test.vtt")
	check(err)
	defer fv.Close()

	r, _ := NewReader(fs)

	for {
		bs := make([]byte, bufSize)
		bv := make([]byte, bufSize)
		ns, _ := r.Read(bs)
		if ns == 0 {
			break
		}
		nv, _ := fv.Read(bv)
		if nv == 0 {
			break
		}
		if bytes.Compare(bs[ns:], bv[nv:]) != 0 {
			t.Errorf("failed real file conversion size %d", bufSize)
		}
	}

}

func TestReader(t *testing.T) {
	testReader(t, 1)
	testReader(t, 64)
	testReader(t, 512)
	testReader(t, 32768)
}
