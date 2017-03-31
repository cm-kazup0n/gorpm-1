package rpmlib

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"io"
	"os"
)

const (
	Null        = 0
	Char        = 1
	Int8        = 2
	Int16       = 3
	Int32       = 4
	Int64       = 5
	String      = 6
	Binary      = 7
	StringArray = 8
	I18nString  = 9
)

const (
	SectionHeaderMagicSize    = 3
	SectionHeaderReservedSize = 4
)
var SectionHeaderMagic []byte = []byte{0x8e, 0xad, 0xe8}

type SectionHeaderIndex struct {
	Tag    int32
	Type   int32
	Offset int32
	Count  int32
}

type SectionHeader struct {
	version int8
	nindex  int32
	hsize   int32
	// reserved 4 bytes
	indexes []SectionHeaderIndex
}

type Section struct {
	magic  []byte
	header *SectionHeader
	store  []byte
}

func readSectionHeader(file *os.File) (header *SectionHeader, err error) {
	header = new(SectionHeader)
	err = binary.Read(file, binary.BigEndian, &header.version)
	if err != nil {
		return
	}

	file.Seek(SectionHeaderReservedSize, os.SEEK_CUR)

	binary.Read(file, binary.BigEndian, &header.nindex)
	if err != nil {
		return
	}

	binary.Read(file, binary.BigEndian, &header.hsize)
	if err != nil {
		return
	}

	header.indexes = make([]SectionHeaderIndex, header.nindex)

	for i, _ := range header.indexes {
		err = binary.Read(file, binary.BigEndian, &header.indexes[i].Tag)
		if err != nil {
			break
		}

		err = binary.Read(file, binary.BigEndian, &header.indexes[i].Type)
		if err != nil {
			break
		}

		err = binary.Read(file, binary.BigEndian, &header.indexes[i].Offset)
		if err != nil {
			break
		}

		err = binary.Read(file, binary.BigEndian, &header.indexes[i].Count)
		if err != nil {
			break
		}
	}

	return
}

func scanSection(file *os.File) (section *Section, err error) {
	section = new(Section)
	section.magic = make([]byte, SectionHeaderMagicSize)

	_, err = file.Read(section.magic)
	if err != nil {
		if err == io.EOF {
			err = fmt.Errorf("Reached EOF before reading a section completed")
		}
		return
	}

	section.header, err = readSectionHeader(file)
	if err != nil {
		if err == io.EOF {
			err = fmt.Errorf("Reached EOF before reading a section completed")
		}
		return
	}

	section.store = make([]byte, section.header.hsize)
	_, err = file.Read(section.store)
	if err != nil {
		return
	}

	err = section.validate()

	return
}

func (section *Section) validate() (err error) {

	// Check magic numbers
	for i, b := range SectionHeaderMagic {
		if b != section.magic[i] {
			return fmt.Errorf("SectionHeader section magic number is invalid")
		}
	}

	var indexes []SectionHeaderIndex
	for i := 0; i < int(section.header.nindex); i++ {
		indexes = append(indexes, section.header.indexes[i])

		for j := len(indexes) - 1; j > 0; j-- {
			if indexes[j-1].Offset > indexes[j].Offset {
				tmp := indexes[j]
				indexes[j] = indexes[j-1]
				indexes[j-1] = tmp
			}
		}
	}

	// TODO
	// Calculate offset and count to check overrun

	return
}

func (section *Section) HasStore(tag int32) (found bool) {
	found = false

	for _, index := range section.header.indexes {
		for index.Tag == tag {
			found = true
			break
		}
	}
	return
}

func (section *Section) GetInt16(tag int32) (value int16, err error) {
	store, _, err := section.GetStore(tag)

	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(store)
	err = binary.Read(buffer, binary.BigEndian, &value)
	
	return
}

func (section *Section) GetInt16Array(tag int32) (value_list []int16, err error) {
	store, _, err := section.GetStore(tag)

	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(store)
	for {
		var value int16
		err := binary.Read(buffer, binary.BigEndian, &value)
		if err != nil {
			break
		}
		value_list = append(value_list, value)	
	}
	return
	
}

func (section *Section) GetInt32(tag int32) (value int32, err error) {
	store, _, err := section.GetStore(tag)

	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(store)
	err = binary.Read(buffer, binary.BigEndian, &value)
	
	return
}

func (section *Section) GetInt32Array(tag int32) (value_list []int32, err error) {
	store, _, err := section.GetStore(tag)

	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(store)
	for {
		var value int32
		err := binary.Read(buffer, binary.BigEndian, &value)
		if err != nil {
			break
		}
		value_list = append(value_list, value)	
	}
	return
	
}

func (section *Section) GetString(tag int32) (value string, err error) {
	store, _, err := section.GetStore(tag)

	if err != nil {
		return
	}

	value = string(store)

	return
}

func (section *Section) GetStringArray(tag int32) (value_list []string, err error) {
	store, _, err := section.GetStore(tag)

	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(store)

	for {
		// byte arrays are separated by NULL byte(0). Read until NULL byte
		a, err := buffer.ReadBytes(0)
		if err != nil {
			break
		}

		// Last index has a NULL byte, Ignore it.
		value_list = append(value_list, string(a[:len(a) - 1]))
	}

	return
}

func (section *Section) GetStore(tag int32) (store []byte, size int32, err error) {
	found := false

	var offset, count, datatype int32
	for _, index := range section.header.indexes {
		for index.Tag == tag {
			found = true
			offset = index.Offset
			count = index.Count
			datatype = index.Type
			break
		}
	}

	if !found {
		return nil, -1, fmt.Errorf("Cannot find store for tag %d", tag)
	}

	switch datatype {
	case Null:
		err = fmt.Errorf("Null field founded.")
		break
	case Char, Int8:
		size = 1
		store = section.store[offset : offset+(count*size)]
		break
	case Int16:
		size = 2
		store = section.store[offset : offset+(count*size)]
		break
	case Int32:
		size = 4
		store = section.store[offset : offset+(count*size)]
		break
	case Int64:
		size = 0
		err = fmt.Errorf("Unsupported type Int64 data type found.")
		break
	case String:
		size = 0
		for _, b := range section.store[offset:] {
			if b == 0 {
				break
			}
			size++
		}
		store = section.store[offset : offset+size]
		break
	case Binary:
		size = 1
		store = section.store[offset : offset+(count*size)]
		break
	case StringArray:
		size = 0
		for i := int32(0); i < count; i++ {
			for _, b := range section.store[offset+size:] {
				size++
				if b == 0 {
					break
				}
			}
		}
		store = section.store[offset : offset+size]
		break
	case I18nString:
		size = 0
		for _, b := range section.store[offset:] {
			if b == 0 {
				break
			}
			size++
		}
		store = section.store[offset : offset+size]
		break
		break
	default:
		err = fmt.Errorf("Unknwon data type %x", datatype)
		break
	}

	return
}
