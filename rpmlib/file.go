package rpmlib

import (
	"encoding/hex"
	"fmt"
	"github.com/pombredanne/gorpm-1/cpio"
	"io"
	"os"
)

type PackageFile struct {
	Lead      *Lead
	Signature *Signature
	Header    *Header
	Payload   *Payload
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

	return
}

type VerifyResult struct {
	Path     string
	FileType int32
	Size     error
	Mode     error
	Checksum error
	MTime    error
}

func (pkg *PackageFile) Verify() (results []VerifyResult, err error) {
	reader := cpio.NewCPIOReader(pkg.Payload.Cpio())

	var archives []*cpio.File

	for {
		file, read_err := reader.GetFile()
		if read_err == io.EOF {
			break
		}

		if read_err != nil {
			err = read_err
			break
		}

		archives = append(archives, file)
	}

	files, err := pkg.Header.Files()
	if err != nil {
		return
	}

	var file_list_on_header []FileMeta
	for _, f := range files {
		if f.Size != 0 {
			file_list_on_header = append(file_list_on_header, f)
		}
	}

	if len(file_list_on_header) != len(archives) {
		err = fmt.Errorf("Number of files are different between header and archive")
		return
	}

	for i, f_h := range file_list_on_header {
		var result VerifyResult

		result.Path = f_h.Path
		result.FileType = f_h.Flag

		if archives[i].Metadata.Filesize != 0 {
			if f_h.Size != int32(archives[i].Metadata.Filesize) {
				result.Size = fmt.Errorf("S: h=%d != a=%d", f_h.Size, archives[i].Metadata.Filesize)
			}
		}

		if f_h.Mode != int16(archives[i].Metadata.Mode) {
			result.Mode = fmt.Errorf("M: h=%x != a=%x", f_h.Mode, archives[i].Metadata.Mode)
			fmt.Println(result.Mode)
		}

		if len(f_h.MD5) > 0 {
			f_checksum, err := hex.DecodeString(f_h.MD5)
			if err != nil {
				return nil, err
			}
			// Why RPMTAG_FILEMD5 stores a sha256 checksum...? Accoding to documents,
			// it stores md5 a checksum....
			if len(f_checksum) == 16 {
				for j, b := range archives[i].MD5() {
					if b != f_checksum[j] {
						result.Checksum = fmt.Errorf("MD5 checksum is invalid")
						break
					}
				}

			} else if len(f_checksum) == 32 {
				for j, b := range archives[i].SHA256() {
					if b != f_checksum[j] {
						result.Checksum = fmt.Errorf("SHA256 checksum is invalid")
						break
					}
				}
			} else {
				// OK, i give up....
			}
		}

		if f_h.Time != int32(archives[i].Metadata.Mtime) {
			result.MTime = fmt.Errorf("Mtime is not match")
		}

		results = append(results, result)
	}

	return
}
