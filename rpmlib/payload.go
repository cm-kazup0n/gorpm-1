package rpmlib

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"github.com/xi2/xz"
)

type Payload struct {
	uncompressed []byte
}


func testread(file *os.File) {
	buffer := make([]byte, 10)
	file.Read(buffer)
	for i, b := range buffer {
		fmt.Printf("%d %x\n", i, b)
	}
}

func getDecompressor(name string, file *os.File) (rd io.Reader, err error) {
	switch name {
	case "xz":
		rd, err = xz.NewReader(file, 0)
		break
	case "gzip":
		rd, err = gzip.NewReader(file)
		break
	default:
		err = fmt.Errorf("Unkown compressor name %s", name)
	}
	return
}

func ScanPayload(file *os.File, comparessor string) (payload *Payload, err error) {

	decompressor, err := getDecompressor(comparessor, file)
	if err != nil {
		return
	}

	payload = new(Payload)
	payload.uncompressed, err = ioutil.ReadAll(decompressor)

	return
}

func (payload *Payload) Cpio() (cpio []byte) {
	cpio = payload.uncompressed

	return
}