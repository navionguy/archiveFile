package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

const tmsg = "The quick brown fox jumped over the lazy dog."

func Test_pumper(t *testing.T) {
	src := ioutil.NopCloser(strings.NewReader(tmsg))

	var p bpump
	p.buf = new(bytes.Buffer)
	cb, err := p.load(src)

	if err != nil {
		t.Errorf("bpump returned error %s, byte count %d/n", err.Error(), cb)
	}

	// okay he loaded, I should be able to successfully
	// unzip it back into the same string

	rd, err := gzip.NewReader(p.buf)

	if err != nil {
		t.Errorf("couldn't put gzip reader on output, error %s/n", err.Error())
	}

	defer rd.Close()

	// read it out

	bf2 := make([]byte, cb)
	_, err = rd.Read(bf2)

	if (err != nil) && (err != io.EOF) {
		t.Errorf("Test_pumper decompress got err = %s/n", err.Error())
	}

	if bytes.Compare([]byte(tmsg), bf2) != 0 {
		t.Errorf("pumper returned %s, expected %s/n", bf2, tmsg)
	}

}
