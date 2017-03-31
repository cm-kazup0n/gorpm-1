package main

import (
	"flag"
	"fmt"
	"github.com/necomeshi/gorpm/rpmlib"
	"os"
)

func PrintPackageInformation(file *os.File) (err error) {

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

	date, err := pkg.Header.BuildDate()
	if err == nil {
		fmt.Printf("BuildDate:  %s\n", date.String())
	} else {
		fmt.Printf("BuildDate:  Unknown\n")
	}
	srpm, err := pkg.Header.SourceRpm()
	if err == nil {
		fmt.Printf("Source RPM: %s\n", srpm)
	}

	fmt.Printf("Summary:    %s\n", pkg.Header.Summary())
	fmt.Printf("Description:\n %s\n", pkg.Header.Description())

	return
}

func PrintPackagedFiles(file *os.File) (err error) {
	pkg, err := rpmlib.ReadPackageFile(file)
	if err != nil {
		return
	}

	file_list, err := pkg.Header.Files()
	if err != nil {
		return
	}

	for _, f := range file_list {
		fmt.Println(f.Path)
	}

	return
}

func PrintPackageFileOf(file *os.File, filetype int32) (err error) {
	pkg, err := rpmlib.ReadPackageFile(file)
	if err != nil {
		return
	}

	file_list, err := pkg.Header.Files()
	if err != nil {
		return
	}

	for _, f := range file_list {
		if f.Flag&filetype != 0 {
			fmt.Println(f.Path)
		}
	}

	return
}

func PrintPackageChangelog(file *os.File) (err error) {
	pkg, err := rpmlib.ReadPackageFile(file)
	if err != nil {
		return
	}

	changelogs, err := pkg.Header.Changelog()
	if err != nil {
		return
	}

	for _, log := range changelogs {
		fmt.Printf("* %s %s\n", log.Date, log.Name)
		fmt.Printf("- %s\n\n", log.Text)
	}

	return
}

func VerifyPackage(file *os.File) (err error) {
	pkg, err := rpmlib.ReadPackageFile(file)
	if err != nil {
		return
	}

	results, err := pkg.Verify()

	if err != nil {
		return
	}

	for _, r := range results {
		if r.Size == nil {
			fmt.Print(".")
		} else {
			fmt.Print("S")
		}

		if r.Mode == nil {
			fmt.Print(".")
		} else {
			fmt.Print("M")
		}

		if r.Checksum == nil {
			fmt.Print(".")
		} else {
			fmt.Print("5")
		}

		if r.MTime == nil {
			fmt.Print(".")
		} else {
			fmt.Print("T")
		}
		fmt.Print(" ")

		switch r.FileType {
		case rpmlib.RPMFILE_CONFIG:
			fmt.Print("c")
			break
		case rpmlib.RPMFILE_DOC:
			fmt.Print("d")
			break
		case rpmlib.RPMFILE_GHOST:
			fmt.Print("g")
			break
		case rpmlib.RPMFILE_LICENSE:
			fmt.Print("l")
			break
		case rpmlib.RPMFILE_README:
			fmt.Print("r")
			break
		default:
			fmt.Print(" ")
			break
		}

		fmt.Print(" ")
		fmt.Println(r.Path)
	}

	return
}

type Option struct {
	ShowInfoMode       bool
	ShowFileMode       bool
	ShowConfigFileMode bool
	ShowDocFileMode    bool
	ShowChangelogMode  bool
	VerificationMode   bool
	//	CheckSignatureMode bool
}

func addOption(option *Option) {
	flag.BoolVar(&option.ShowInfoMode, "i", false, "Show package inforamtion.")
	flag.BoolVar(&option.ShowFileMode, "l", false, "Show files included package.")
	flag.BoolVar(&option.ShowConfigFileMode, "c", false, "Show config files included package.")
	flag.BoolVar(&option.ShowDocFileMode, "d", false, "Show doc files included package.")
	flag.BoolVar(&option.ShowChangelogMode, "changelog", false, "Show changelog.")
	flag.BoolVar(&option.VerificationMode, "V", false,
		"Verify file's size, checksum, permission and type. user and group are not verified.")
	//flag.BoolVar(&option.CheckSignatureMode, "checksig", false, "Check all digests and signatures"
}

func main() {
	var option Option

	addOption(&option)

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Package file not specified\n")
		os.Exit(1)
	}

	for _, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		if option.ShowInfoMode {
			err = PrintPackageInformation(file)
		} else if option.ShowFileMode {
			err = PrintPackagedFiles(file)
		} else if option.ShowConfigFileMode {
			err = PrintPackageFileOf(file, rpmlib.RPMFILE_CONFIG)
		} else if option.ShowDocFileMode {
			err = PrintPackageFileOf(file, rpmlib.RPMFILE_DOC)
		} else if option.ShowChangelogMode {
			err = PrintPackageChangelog(file)
		} else if option.VerificationMode {
			err = VerifyPackage(file)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		file.Close()
	}

	os.Exit(0)
}
