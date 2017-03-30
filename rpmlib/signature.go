package rpmlib

import (
	"bytes"
	"encoding/binary"
	"os"
	"fmt"
)

const (
	RPMSIGTAG_DSA         = 267
	RPMSIGTAG_RSA         = 268
	RPMSIGTAG_SHA1        = 269
	RPMSIGTAG_SIZE        = 1000
	RPMSIGTAG_PGP         = 1002
	RPMSIGTAG_MD5         = 1004
	RPMSIGTAG_GPG         = 1005
	RPMSIGTAG_PAYLOADSIZE = 1007
	RPMSIGTAG_SAH1HEADER  = 1010
)

type Signature struct {
	Section
}

//
// Required
//
func (sig *Signature) Size() (size int32, err error) {

	store, _, err := sig.GetStore(RPMSIGTAG_SIZE)
	if err != nil {
		return
	}
	err = binary.Read(bytes.NewReader(store), binary.BigEndian, &size)

	return
}

func (sig *Signature) MD5() (bin []byte, err error) {

	bin, nsize, err := sig.GetStore(RPMSIGTAG_MD5)

	if nsize != 16 {
		return nil, fmt.Errorf("Less size for required field 'MD5' %d", nsize)	
	}

	if err != nil {
		return
	}

	return
}

func ScanSignature(file *os.File) (signature *Signature, err error) {

	section, err := scanSection(file)
	if err != nil {
		return
	}

	signature = new(Signature)
	signature.Section = *section

	return
}

//
// Optional
//
func (sig *Signature) HasPayloadSize() (hasPayloadSize bool) {
	return sig.HasStore(RPMSIGTAG_PAYLOADSIZE)
}

func (sig *Signature) HasSAH1() (hasSAH1 bool) {
	return sig.HasStore(RPMSIGTAG_SAH1HEADER)
}

func (sig *Signature) PayloadSize() (size int32, err error) {

	store, _, err := sig.GetStore(RPMSIGTAG_PAYLOADSIZE)
	if err != nil {
		return
	}
	err = binary.Read(bytes.NewReader(store), binary.BigEndian, &size)

	return
}

func (sig *Signature) SAH1() (checksum []byte, err error) {

	checksum, _, err = sig.GetStore(RPMSIGTAG_PAYLOADSIZE)
	if err != nil {
		return
	}

	return
}