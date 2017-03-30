package rpmlib

import (
	"os"
//	"fmt"
)

type PackageFile struct {
	Lead      *Lead
	Signature *Signature
	Header    *Header
	Payload *Payload
}

func ReadPackageFile(file *os.File) (pkg *PackageFile, err error) {
	pkg = new(PackageFile)

	pkg.Lead, err = ScanLead(file)
	if err != nil {
		return
	}

	pkg.Signature, err = ScanSignature(file)

	if err != nil {
		return
	}

	// Data store size is always 8 byte boundary
	boundary := pkg.Signature.header.hsize % 8
	if boundary != 0 {
		file.Seek(int64(boundary), os.SEEK_CUR)
	}

	pkg.Header, err = ScanHeader(file)

	if err != nil {
		return
	}

	compressor := pkg.Header.PayloadCompressor()
	pkg.Payload, err = ScanPayload(file, compressor)

	if err != nil {
		return
	}

//	for _, index := range pkg.Header.header.indexes {
//		fmt.Printf("tag = %05d type = %05d count =%d\n", index.Tag, index.Type, index.Count)
//	}

	return
}