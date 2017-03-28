package main

import (
	"flag"
	"fmt"
	"github.com/necomeshi/gorpm/rpmlib"
	"os"
)

func PrintPackageInformation(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	pkg, err := rpmlib.ReadPackageFile(file)
	if err != nil {
		return
	}

	fmt.Printf("Name:       %s\n", pkg.Header.Name())
	fmt.Printf("Version:    %s\n", pkg.Header.Version())
	fmt.Printf("Release:    %s\n", pkg.Header.Release())
	fmt.Printf("Group:      %s\n", pkg.Header.Group())
	fmt.Printf("Size:       %d\n", pkg.Header.Size())
	fmt.Printf("Licence:    %s\n", pkg.Header.Licence())
	fmt.Printf("BuildDate:  %s\n", pkg.Header.BuildDate().String())
	fmt.Printf("Source RPM: %s\n", pkg.Header.SourceRpm())
	fmt.Printf("Summary:    %s\n", pkg.Header.Summary())
	fmt.Printf("Description:\n %s\n", pkg.Header.Description())

	return
}

func PrintPackagedFiles(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	pkg, err := rpmlib.ReadPackageFile(file)
	if err != nil {
		return
	}

	filenames, err := pkg.Header.FileList()
	if err != nil {
		return
	}

	for _, name := range filenames {
		fmt.Println(name)
	}

	return
}

func PrintPackageChangelog(filename string) (err error) {

	return
}

type Option struct {
	ShowInfoMode bool
	ShowFileMode bool
	ShowChangelogMode bool
}

func addOption(option *Option) {
	flag.BoolVar(&option.ShowInfoMode, "i", false, "Show package inforamtion.")
	flag.BoolVar(&option.ShowFileMode, "l", false, "Show package files.")
	flag.BoolVar(&option.ShowFileMode, "c", false, "Show changelog.")
}

func main() {
	var option Option

	addOption(&option)

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Package file not specified\n")
		os.Exit(1)
	}

	if option.ShowInfoMode {
		for _, filename := range flag.Args() {
			err := PrintPackageInformation(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		}
	} else if option.ShowFileMode {
		for _, filename := range flag.Args() {
			err := PrintPackagedFiles(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		}
	} else if option.ShowChangelogMode {
		
	}

	os.Exit(0)
}
