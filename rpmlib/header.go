package rpmlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
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

func (header *Header) PayloadCompressor() (name string) {
	store, _, _ := header.Section.GetStore(RPMTAG_PAYLOADCOMPRESSOR)

	name = string(store)
	
	return
}

func (header *Header) FileFlags() (flags []int32) {
	store, _, _ := header.Section.GetStore(RPMTAG_FILEFLAGS)	

	reader := bytes.NewReader(store)
	
	for {
		var flag int32
		err := binary.Read(reader, binary.BigEndian, &flag)
		if err != nil {
			break
		}

		flags = append(flags, flag)
	}
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
	if header.Section.HasStore(RPMTAG_BASENAMES) &&
		header.Section.HasStore(RPMTAG_DIRNAMES) &&
		header.Section.HasStore(RPMTAG_DIRINDEXES) {

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
			filenames = append(filenames, dirnames[dirindexes[i]]+basename)
		}
	} else if header.Section.HasStore(RPMTAG_OLDFILENAMES) {
		store, _, _ := header.Section.GetStore(RPMTAG_OLDFILENAMES)
		buffer := bytes.NewBuffer(store)

		for {
			s, err := buffer.ReadString(0)
			if err != nil {
				break
			}
			filenames = append(filenames, s)
		}
	} else {
		err = fmt.Errorf("File list data not presented")
	}

	return
}

func (header *Header) Changelog() (logs []Changelog, err error) {
	if !header.Section.HasStore(RPMTAG_CHANGELOGNAME) ||
		!header.Section.HasStore(RPMTAG_CHANGELOGTEXT) ||
		!header.Section.HasStore(RPMTAG_CHANGELOGTIME) {

		return nil, fmt.Errorf("No changelog data found")
	}

	// Already checked it exists.
	store, _, _ := header.Section.GetStore(RPMTAG_CHANGELOGNAME)

	buffer := bytes.NewBuffer(store)

	var lognames []string
	for {
		s, err := buffer.ReadString(0)
		if err != nil {
			break
		}
		lognames = append(lognames, s)
	}

	store, _, _ = header.Section.GetStore(RPMTAG_CHANGELOGTEXT)

	buffer = bytes.NewBuffer(store)

	var logtexts []string
	for {
		s, err := buffer.ReadString(0)
		if err != nil {
			break
		}
		logtexts = append(logtexts, s)
	}

	store, _, _ = header.Section.GetStore(RPMTAG_CHANGELOGTIME)

	buffer = bytes.NewBuffer(store)

	var logtimes []time.Time
	for {
		var unixtime int32
		err := binary.Read(buffer, binary.BigEndian, &unixtime)
		if err != nil {
			break
		}
		logtimes = append(logtimes, time.Unix(int64(unixtime), 0))
	}

	if len(lognames) == len(logtexts) && len(lognames) == len(logtimes) {
		for i, _ := range lognames {
			logs = append(logs, Changelog{lognames[i], logtexts[i], logtimes[i]})
		}
	} else {
		return nil, fmt.Errorf("Changelog's name, text, time array size are different")
	}

	return
}

type Changelog struct {
	Name string
	Text string
	Date time.Time
}
