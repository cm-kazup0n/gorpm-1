package rpmlib

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"encoding/binary"
	"xi2.org/x/xz"
)

const (
	CPIO_HEADER_SIZE = 110
	CPIO_HEADER_MAGIC_SIZE = 6
)

type Payload struct {
	uncompressed []byte
}

type FileMeta struct {
	name	 string
	ino       uint64
	mode      uint64
	uid       uint64
	gid       uint64
	nlink     uint64
	mtime     uint64
	filesize  uint64
	devmajor  uint64
	devminor  uint64
	rdevmajor uint64
	rdevminor uint64
	namesize  uint64
	checksum  uint64
}

func NewFileMeta(cpioheader []byte, name string) (meta FileMeta, err error) {
	meta.name = name

	reader := bytes.NewReader(cpioheader)

	magic := make([]byte, CPIO_HEADER_MAGIC_SIZE)

	reader.Read(magic)


	var table []*uint64 = []*uint64{
		&meta.ino, &meta.mode, &meta.uid, &meta.gid, &meta.nlink,
		&meta.mtime, &meta.filesize, &meta.devmajor, &meta.devminor,
		&meta.rdevmajor, &meta.rdevminor, &meta.namesize, &meta.checksum,
	}

	for _, ptr := range table {
		*ptr, err = binary.ReadUvarint(reader)
		if err != nil {
			break
		}
	}

	return
}

type File struct {
	FileMeta
	data []byte
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

//func (payload *Payload) Files() (files []File, err error) {
//
//	reader := bytes.NewBuffer(payload.uncompressed)
//	cpioheader := make([]byte, CPIO_HEADER_SIZE)
//
//	for {
//
//		file := new(File)
//
//		n, err := reader.Read(cpioheader)
//		if err != nil {
//			break
//		}
//
//		magic := cpioheader[0:CPIO_HEADER_MAGIC_SIZE]
//		fmt.Println(magic)
//	
//
//		name, err := reader.ReadBytes(0)
//		if err != nil {
//			break
//		}
//
//		file.FileMeta, err = NewFileMeta(cpioheader, string(name))
//		if err != nil {
//			break
//		}
//
//		skip := (n + len(name) + 1) % 4
//		if skip != 0 {
//			reader.Next(skip)
//		}
//
//
//		file.data = make([]byte, file.FileMeta.filesize)
//		n, err = reader.Read(file.data)
//
//		if err != nil {
//			break
//		}
//
//		skip = n % 4
//		if skip != 0 {
//			reader.Next(skip)
//		}
//	}
//
//	return
//}
