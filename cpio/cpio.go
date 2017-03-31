package cpio

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"strconv"
)

const (
	CPIO_NEW_HEADER_SIZE       = 110
	CPIO_NEW_HEADER_MAGIC_SIZE = 6
	CPIO_NEW_FIELD_SIZE        = 8
)

type CPIOType int

const (
	CPIO_NEW_ASCII CPIOType = iota
	CPIO_NEW_CRC
)

type Meta struct {
	Type      CPIOType
	Size      int
	Ino       uint64
	Mode      uint64
	Uid       uint64
	Gid       uint64
	Nlink     uint64
	Mtime     uint64
	Filesize  uint64
	Devmajor  uint64
	Devminor  uint64
	Rdevmajor uint64
	Rdevminor uint64
	Namesize  uint64
	Checksum  uint64
}

func NewMetadata(rd *bytes.Reader) (m *Meta, err error) {
	magic_bytes := make([]byte, CPIO_NEW_HEADER_MAGIC_SIZE)

	_, err = rd.Read(magic_bytes)
	if err != nil {
		return
	}

	meta := new(Meta)

	//
	// NOTE
	// Old CPIO format is not supported currently
	//
	switch string(magic_bytes) {
	case "070701":
		meta.Type = CPIO_NEW_ASCII
		meta.Size = CPIO_NEW_HEADER_SIZE
		break
	case "070702":
		meta.Type = CPIO_NEW_CRC
		meta.Size = CPIO_NEW_HEADER_SIZE
		break
	default:
		err = fmt.Errorf("CPIO magic number invalid or not supported")
		return
	}

	field := make([]byte, CPIO_NEW_FIELD_SIZE)
	var table []*uint64 = []*uint64{
		&meta.Ino, &meta.Mode, &meta.Uid, &meta.Gid, &meta.Nlink,
		&meta.Mtime, &meta.Filesize, &meta.Devmajor, &meta.Devminor,
		&meta.Rdevmajor, &meta.Rdevminor, &meta.Namesize, &meta.Checksum,
	}

	for _, pv := range table {
		_, err = rd.Read(field)
		if err != nil {
			return
		}
		*pv, err = strconv.ParseUint(string(field), 16, 64)
		if err != nil {
			return
		}
	}

	return meta, err
}

type File struct {
	Metadata *Meta
	Name     string
	data     []byte
}

func (f *File) Write(w io.Writer) (n int, err error) {
	return w.Write(f.data)
}

func (f *File) MD5() [md5.Size]byte {
	return md5.Sum(f.data)
}

func (f *File) SHA256() [sha256.Size]byte {
	return sha256.Sum256(f.data)
}

type CPIOReader struct {
	reader *bytes.Reader
}

func NewCPIOReader(cpiodata []byte) (rd *CPIOReader) {
	rd = new(CPIOReader)
	rd.reader = bytes.NewReader(cpiodata)

	return
}

func (rd *CPIOReader) GetFile() (file *File, err error) {
	file = new(File)
	file.Metadata, err = NewMetadata(rd.reader)
	if err != nil {
		return
	}

	var name []byte
	for {
		b, rerr := rd.reader.ReadByte()
		if rerr != nil {
			return
		}
		if b == 0 {
			break
		}
		name = append(name, b)
	}
	file.Name = string(name)

	if file.Name == "TRAILER!!!" {
		return nil, io.EOF
	}

	// last NULL byte is already read. So plus 1
	for total := file.Metadata.Size + len(name) + 1; total%4 != 0; total++ {
		_, err = rd.reader.ReadByte()
		if err != nil {
			return
		}
	}

	file.data = make([]byte, file.Metadata.Filesize)

	n, err := rd.reader.Read(file.data)
	if err != nil {
		return
	}

	if uint64(n) != file.Metadata.Filesize {
		err = fmt.Errorf("Cannot read file data. less than %d", file.Metadata.Filesize)
		return
	}

	for ; n%4 != 0; n++ {
		_, err = rd.reader.ReadByte()
		if err != nil {
			return
		}
	}

	return
}
