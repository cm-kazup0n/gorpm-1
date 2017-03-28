package rpmlib

import (
	"os"
	"fmt"
	"bytes"
	"encoding/binary"
	"time"
)

var HeaderRequiredField []int32 = []int32{
	RPMTAG_HEADER18NTABLE,
	RPMTAG_NAME,
	RPMTAG_VERSION,
	RPMTAG_RELEASE,
	RPMTAG_SUMMARY,
	RPMTAG_DESCRIPTION,
	RPMTAG_SIZE,
	RPMTAG_LICENCE,
	RPMTAG_GROUP,
	RPMTAG_OS,
	RPMTAG_ARCH,
	RPMTAG_PAYLOADFORMAT,
	RPMTAG_PAYLOADCOMPRESSOR,
	RPMTAG_PAYLOAD_FLAGS,
}

type Header struct {
	Section
}

func ScanHeader(file *os.File) (header *Header, err error) {
	err = SkipSignature(file)
	if err != nil {
		return
	}

	section, err := scanSection(file)
	if err != nil {
		return
	}

	header = new(Header)
	header.Section = *section

	for _, tag := range HeaderRequiredField {
		if !header.Section.HasStore(tag) {
			err = fmt.Errorf("Cannot find required field tag=%d", tag)
			break
		}
	}

	return
}

//
// Required Fields
// These field shall present and already checked above
// so no error will happen
//
func (header *Header) Name() (name string) {
	store, _, _ := header.Section.GetStore(RPMTAG_NAME)

	name = string(store)

	return
}

func (header *Header) Version() (version string) {
	store, _, _ := header.Section.GetStore(RPMTAG_VERSION)

	version = string(store)

	return
}

func (header *Header) Release() (release string) {
	store, _, _ := header.Section.GetStore(RPMTAG_RELEASE)

	release = string(store)

	return
}

func (header *Header) Group() (group string) {
	store, _, _ := header.Section.GetStore(RPMTAG_GROUP)

	group = string(store)

	return
}

func (header *Header) Size() (size int32) {
	store, _, _ := header.Section.GetStore(RPMTAG_SIZE)

	binary.Read(bytes.NewReader(store), binary.BigEndian, &size)

	return
}

func (header *Header) Summary() (summary string) {
	store, _, _ := header.Section.GetStore(RPMTAG_SUMMARY)

	summary = string(store)

	return
}

func (header *Header) Description() (description string) {
	store, _, _ := header.Section.GetStore(RPMTAG_DESCRIPTION)

	description = string(store)

	return
}

func (header *Header) Licence() (licence string) {
	store, _, _ := header.Section.GetStore(RPMTAG_LICENCE)

	licence = string(store)

	return
}

func (header *Header) SourceRpm() (name string) {
	store, _, _ := header.Section.GetStore(RPMTAG_SOURCERPM)

	name = string(store)

	return
}

func (header *Header) BuildDate() (buildtime time.Time) {
	store, _, _ := header.Section.GetStore(RPMTAG_BUILDTIME)

	var t int32
	binary.Read(bytes.NewReader(store), binary.BigEndian, &t)
	
	buildtime = time.Unix(int64(t), 0)

	return
}

func (header *Header) FileList() (filenames []string, err error) {
	if !header.Section.HasStore(RPMTAG_BASENAMES) ||
		!header.Section.HasStore(RPMTAG_DIRNAMES) ||
		!header.Section.HasStore(RPMTAG_DIRINDEXES) {
		return
	}

	// Get filename list(basename)
	store, _, _ := header.Section.GetStore(RPMTAG_BASENAMES)
	buffer := bytes.NewBuffer(store)

	var basenames []string
	for {
		s, err := buffer.ReadString(0)
		if err != nil {
			break
		}
		basenames = append(basenames, s)
	}

	store, _, _ = header.Section.GetStore(RPMTAG_DIRNAMES)
	buffer = bytes.NewBuffer(store)

	var dirnames []string
	for {
		s, err := buffer.ReadString(0)
		if err != nil {
			break
		}
		dirnames = append(dirnames, s)
	}

	var dirindexes []int32
	store, _, _ = header.Section.GetStore(RPMTAG_DIRINDEXES)
	reader := bytes.NewReader(store)
	for {
		var index int32
		readerr := binary.Read(reader, binary.BigEndian, &index)
		if readerr != nil {
			break
		}	
		dirindexes = append(dirindexes, index)
	}

	if len(dirindexes) != len(basenames) {
		return nil, fmt.Errorf("directory indexes length differente from length of basenames")
	}

	for i, basename := range basenames {
		filenames = append(filenames, dirnames[dirindexes[i]] + basename)
	}
	
	return
}