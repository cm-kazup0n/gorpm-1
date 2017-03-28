package rpmlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	SupportedMajorVersion = 3
	SupportedMinorVersion = 0
)

const (
	BinaryPackageFileType uint16 = 0x0000
	SourcePackageFileType uint16 = 0x0001
)

const (
	LeadSize              int64 = 96
	LeadMagicSize               = 4
	LeadMajorSize               = 1
	LeadMinorSize               = 1
	LeadRpmTypeSize             = 2
	LeadArchNumSize             = 2
	LeadNameSize                = 66
	LeadOSNumSize               = 2
	LeadSignatureTypeSize       = 2
	// reserved 16 bytes
)

var LeadMagic = []byte{0xed, 0xab, 0xee, 0xdb}

type Lead struct {
	data []byte
}

func (lead *Lead) Validate() (err error) {

	magic := lead.Magic()

	for i, b := range LeadMagic {
		if magic[i] != b {
			return fmt.Errorf("Lead maigc is invalid")
		}
	}

	// is supported version ?
	if lead.Major() != SupportedMajorVersion || lead.Minor() != SupportedMinorVersion {
		return fmt.Errorf("This rpm format version is not supported.")
	}

	return
}

func (lead *Lead) Magic() (magic []uint8) {
	magic = make([]uint8, LeadMagicSize)
	buffer := bytes.NewBuffer(lead.data)

	binary.Read(buffer, binary.BigEndian, &magic)

	return
}

func (lead *Lead) Major() (major uint8) {
	offset := LeadMagicSize
	buffer := bytes.NewBuffer(lead.data[offset:])
	binary.Read(buffer, binary.BigEndian, &major)

	return
}

func (lead *Lead) Minor() (minor uint8) {
	offset := LeadMagicSize + LeadMajorSize
	buffer := bytes.NewBuffer(lead.data[offset:])
	binary.Read(buffer, binary.BigEndian, &minor)

	return
}

func (lead *Lead) RpmType() (rpmtype uint16) {
	offset := LeadMagicSize + LeadMajorSize + LeadMinorSize
	buffer := bytes.NewBuffer(lead.data[offset:])
	binary.Read(buffer, binary.BigEndian, &rpmtype)

	return
}

func (lead *Lead) ArchtectureNumber() (archnum uint16) {
	offset := LeadMagicSize + LeadMajorSize + LeadMinorSize + LeadRpmTypeSize
	buffer := bytes.NewBuffer(lead.data[offset:])
	binary.Read(buffer, binary.BigEndian, &archnum)

	return
}

func (lead *Lead) Name() (name string) {
	offset := LeadMagicSize + LeadMajorSize + LeadMinorSize + LeadRpmTypeSize +
		LeadArchNumSize

	buffer := bytes.NewBuffer(lead.data[offset:])
	name = buffer.String()

	return
}

func (lead *Lead) OSNumber() (number int16) {
	offset := LeadMagicSize + LeadMajorSize + LeadMinorSize + LeadRpmTypeSize +
		LeadArchNumSize + LeadNameSize

	buffer := bytes.NewBuffer(lead.data[offset:])
	binary.Read(buffer, binary.BigEndian, &number)

	return
}

func (lead *Lead) SignatureType() (typenumber int16) {
	offset := LeadMagicSize + LeadMajorSize + LeadMinorSize + LeadRpmTypeSize +
		LeadArchNumSize + LeadNameSize + LeadOSNumSize

	buffer := bytes.NewBuffer(lead.data[offset:])
	binary.Read(buffer, binary.BigEndian, &typenumber)

	return
}

func ScanLead(file *os.File) (rpmlead *Lead, err error) {

	rpmlead = new(Lead)
	rpmlead.data = make([]byte, LeadSize)

	nsize, err := file.Read(rpmlead.data)

	if err != nil {
		return
	}

	if int64(nsize) != LeadSize {
		return nil, fmt.Errorf("Less data size for RPM Lead")
	}

	err = rpmlead.Validate()

	return
}

func SkipLead(file *os.File) (err error) {
	_, err = file.Seek(LeadSize, os.SEEK_SET)

	return
}
