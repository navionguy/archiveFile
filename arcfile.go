package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

const buffSize = 5000

type bpump struct {
	buf *bytes.Buffer // the buffer where data gets compressed
}

type uploadManager struct {
}

type result struct {
	location string // where the data went
	tb       int64  // bytes stored
}

// ArchiveFile takes a src file, compresses it using gzip and then writes it to the supplied s3Store
func ArchiveFile(src io.ReadCloser, store uploadManager) error {
	// create my pump
	var p bpump
	p.buf = new(bytes.Buffer)

	// load and compress the soure file
	_, err := p.load(src)

	// then go tell the store to upload the data

	res, err := store.upload(p)

	if err != nil {
		return err
	}

	fmt.Printf("file successfully archived at %s/n", res.location)

	return nil
}

// loads the uncompressed data into a compressed file
func (p bpump) load(f io.Reader) (int64, error) {
	// create a gzip writer over my buffer
	gz := gzip.NewWriter(p.buf)
	defer gz.Close()

	// just copy it all
	cb, err := io.Copy(gz, f)

	return cb, err
}

// reads as many bytes as the passed buffer will hold
func (p bpump) Read(data []byte) (cb int, err error) {
	cb, err = p.buf.Read(data)

	return cb, err
}
