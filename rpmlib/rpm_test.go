package rpmlib

import (
    "testing"
	"os"
	"fmt"
)

const (
	SampleRpm1 string = "/Users/yoshimura/Downloads/rpm-4.8.0-55.el6.x86_64.rpm"
)

func testread(file *os.File) {
	buffer := make([]byte, 10)
	file.Read(buffer)
	for i, b := range buffer {
		fmt.Printf("%d %x\n", i, b)
	}
}


func TestXYZ(t *testing.T) {
	file, _ := os.Open(SampleRpm1)

	pkg, err := ReadPackageFile(file)
	if err != nil {
		t.Fatal(err)	
	}

	pkg.Header.Name()
}
