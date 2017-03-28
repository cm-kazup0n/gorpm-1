package rpmlib

import (
	"os"
	"fmt"
)

type PackageFile struct {
	Lead      *Lead
	Signature *Signature
	Header    *Header
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

	pkg.Header, err = ScanHeader(file)

	if err != nil {
		return
	}

	// for debug
	for _, index := range pkg.Header.header.indexes {
		fmt.Printf("tag = %05d type = %05d count =%d\n", index.Tag, index.Type, index.Count)
	}

	return
}