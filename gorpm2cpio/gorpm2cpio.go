package main 

import (
	"fmt"
	"github.com/necomeshi/gorpm/rpmlib"
	"os"
)

func main() {

	if len(os.Args) == 0 {
		fmt.Fprintf(os.Stderr, "No package file specified")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open the file '%s' : %s", os.Args[1], err)
		os.Exit(1)
	}

	pkg, err := rpmlib.ReadPackageFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading file: %s", err)
		os.Exit(1)
	}

	cpio := pkg.Payload.Cpio()

	n, err := os.Stdout.Write(cpio)

	if n != len(cpio) {
		fmt.Fprintf(os.Stderr, "Cannot write all bytes to stdout")
	}

	return
}

