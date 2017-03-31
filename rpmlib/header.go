package rpmlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

const (
	// Header private
	RPMTAG_HEADERSIGNATURES = 62
	RPMTAG_HEADERIMMUTABLE  = 63
	RPMTAG_HEADER18NTABLE   = 100

	RPMTAG_NAME              = 1000
	RPMTAG_VERSION           = 1001
	RPMTAG_RELEASE           = 1002
	RPMTAG_SUMMARY           = 1004
	RPMTAG_DESCRIPTION       = 1005
	RPMTAG_BUILDTIME         = 1006
	RPMTAG_BUILDHOST         = 1007
	RPMTAG_SIZE              = 1009
	RPMTAG_DISTRIBUTION      = 1010
	RPMTAG_VERNDOR           = 1011
	RPMTAG_LICENCE           = 1014
	RPMTAG_PACKAGER          = 1015
	RPMTAG_GROUP             = 1016
	RPMTAG_URL               = 1020
	RPMTAG_OS                = 1021
	RPMTAG_ARCH              = 1022
	RPMTAG_OLDFILENAMES      = 1027
	RPMTAG_FILESIZES         = 1028
	RPMTAG_FILEMODES         = 1030
	RPMTAG_FILERDEVS         = 1033
	RPMTAG_FILEMTIMES        = 1034
	RPMTAG_FILEMD5S          = 1035
	RPMTAG_FILEDIGESTS       = 1035
	RPMTAG_FILELINKTOS       = 1036
	RPMTAG_FILEFLAGS         = 1037
	RPMTAG_FILEUSERNAME      = 1039
	RPMTAG_FILEGROUPNAME     = 1040
	RPMTAG_SOURCERPM         = 1044
	RPMTAG_ARCHIVESIZE       = 1046
	RPMTAG_RPMVERSION        = 1064
	RPMTAG_CHANGELOGTIME     = 1080
	RPMTAG_CHANGELOGNAME     = 1081
	RPMTAG_CHANGELOGTEXT     = 1082
	RPMTAG_COOKIE            = 1094
	RPMTAG_FILEDEVICES       = 1095
	RPMTAG_FILEINODES        = 1096
	RPMTAG_FILELANGS         = 1097
	RPMTAG_DIRINDEXES        = 1116
	RPMTAG_BASENAMES         = 1117
	RPMTAG_DIRNAMES          = 1118
	RPMTAG_DISTURL           = 1123
	RPMTAG_PAYLOADFORMAT     = 1124
	RPMTAG_PAYLOADCOMPRESSOR = 1125
	RPMTAG_PAYLOAD_FLAGS     = 1126
)

const (
	RPMFILE_NONE      = 0
	RPMFILE_CONFIG    = 1
	RPMFILE_DOC       = 1 << 1
	RPMFILE_ICON      = 1 << 2
	RPMFILE_MISSINGOK = 1 << 3
	RPMFILE_NOREPLACE = 1 << 4
	RPMFILE_SPECFILE  = 1 << 5
	RPMFILE_GHOST     = 1 << 6
	RPMFILE_LICENSE   = 1 << 7
	RPMFILE_README    = 1 << 8
	RPMFILE_UNPATCHED = 1 << 9
	RPMFILE_PUBKEY    = 1 << 10
	RPMFILE_POLICY    = 1 << 11
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

type FileMeta struct {
	Name    string
	Path    string
	Size    int32
	Mode    int16
	Device  int32
	Time    int32
	MD5     string
	LinkTo  string
	Flag    int32
	User    string
	Group   string
	RDevice int16
	Inode   int32
	Lang    string
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

func (header *Header) Files() (meta_list []FileMeta, err error) {
	filenames, err := header.FileNames()
	if err != nil {
		return
	}

	size_list, err := header.Section.GetInt32Array(RPMTAG_FILESIZES)
	if err != nil {
		return
	}

	flag_list, err := header.Section.GetInt32Array(RPMTAG_FILEFLAGS)
	if err != nil {
		return
	}

	md5_list, err := header.Section.GetStringArray(RPMTAG_FILEMD5S)
	if err != nil {
		return
	}

	rdevnum_list, err := header.Section.GetInt16Array(RPMTAG_FILERDEVS)
	if err != nil {
		return
	}

	devnum_list, err := header.Section.GetInt32Array(RPMTAG_FILEDEVICES)
	if err != nil {
		return
	}

	linkto_list, err := header.Section.GetStringArray(RPMTAG_FILELINKTOS)
	if err != nil {
		return
	}

	mtime_list, err := header.Section.GetInt32Array(RPMTAG_FILEMTIMES)
	if err != nil {
		return
	}

	mode_list, err := header.Section.GetInt16Array(RPMTAG_FILEMODES)
	if err != nil {
		return
	}

	if len(filenames) != len(size_list) || len(filenames) != len(md5_list) ||
		len(filenames) != len(rdevnum_list) || len(filenames) != len(devnum_list) ||
		len(filenames) != len(linkto_list) || len(filenames) != len(mtime_list) ||
		len(filenames) != len(mode_list) {

		err = fmt.Errorf("Number of file's name, attributes different")
		return
	}

	for i, name := range filenames {
		var meta FileMeta
		meta.Path = name
		meta.Size = size_list[i]
		meta.Flag = flag_list[i]
		meta.MD5 = md5_list[i]
		meta.Device = devnum_list[i]
		meta.RDevice = rdevnum_list[i]
		meta.LinkTo = linkto_list[i]
		meta.Time = mtime_list[i]
		meta.Mode = mode_list[i]

		meta_list = append(meta_list, meta)
	}

	return
}

func (header *Header) FileNames() (filenames []string, err error) {
	if header.Section.HasStore(RPMTAG_BASENAMES) &&
		header.Section.HasStore(RPMTAG_DIRNAMES) &&
		header.Section.HasStore(RPMTAG_DIRINDEXES) {

		basenames, err := header.Section.GetStringArray(RPMTAG_BASENAMES)
		if err != nil {
			return nil, err
		}

		dirnames, err := header.Section.GetStringArray(RPMTAG_DIRNAMES)
		if err != nil {
			return nil, err
		}

		dirindexes, err := header.Section.GetInt32Array(RPMTAG_DIRINDEXES)
		if err != nil {
			return nil, err
		}

		if len(dirindexes) != len(basenames) {
			return nil, fmt.Errorf("directory indexes length differente from length of basenames")
		}

		for i, basename := range basenames {
			filenames = append(filenames, dirnames[dirindexes[i]]+basename)
		}

	} else if header.Section.HasStore(RPMTAG_OLDFILENAMES) {
		filenames, err = header.Section.GetStringArray(RPMTAG_OLDFILENAMES)
	} else {
		err = fmt.Errorf("File list data not presented")
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

//
// Optional or Informational field
//
func (header *Header) SourceRpm() (name string, err error) {
	return header.Section.GetString(RPMTAG_SOURCERPM)
}

func (header *Header) BuildDate() (buildtime time.Time, err error) {

	t, err := header.Section.GetInt32(RPMTAG_BUILDTIME)
	if err != nil {
		return
	}

	buildtime = time.Unix(int64(t), 0)

	return
}

func (header *Header) Changelog() (logs []Changelog, err error) {
	if !header.Section.HasStore(RPMTAG_CHANGELOGNAME) ||
		!header.Section.HasStore(RPMTAG_CHANGELOGTEXT) ||
		!header.Section.HasStore(RPMTAG_CHANGELOGTIME) {

		return nil, fmt.Errorf("No changelog data found")
	}

	lognames, err := header.Section.GetStringArray(RPMTAG_CHANGELOGNAME)
	if err != nil {
		return
	}

	logtexts, err := header.Section.GetStringArray(RPMTAG_CHANGELOGTEXT)
	if err != nil {
		return
	}

	unix_times, err := header.Section.GetInt32Array(RPMTAG_CHANGELOGTIME)
	if err != nil {
		return
	}

	var logtimes []time.Time
	for _, unix_time := range unix_times {
		logtimes = append(logtimes, time.Unix(int64(unix_time), 0))
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
