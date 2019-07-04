package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type fpump struct {
	odata *os.File // file that I'm pretending is a store
}

//
func main() {

	infile := os.Args[1] // first parameter is my input file

	inp, err := os.Open(infile)

	if err != nil {
		fmt.Printf("error %s opening input file/n", err.Error())
		return
	}

	// create a mock store that just writes to a temp file
	var myStore uploadManager

	err = ArchiveFile(inp, myStore)

	if (err != nil) && (err != io.EOF) {
		fmt.Printf("ArchvieFile failed with error, %s/n", err.Error())
	}
}

// upload moves the data to "mock" S3 storage
func (um uploadManager) upload(f io.Reader) (res result, err error) {
	// open a temporary file which we will treat like S3 storage
	tf, err := ioutil.TempFile("", "S3_")

	if err != nil {
		fmt.Printf("creating storage file encountered error %s/n", err.Error())
		return
	}
	defer tf.Close()

	// tell caller where the file went
	res.location = tf.Name()

	// copy everything
	res.tb, err = io.Copy(tf, f)

	if (err != nil) && (err != io.EOF) {
		fmt.Printf("error copying file to storage, %s/n", err.Error())
	}

	return res, err
}
