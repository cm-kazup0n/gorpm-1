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
	fmt.Printf("BuildDate:  %s\n", pkg.Header.BuildDate().String())
	fmt.Printf("Source RPM: %s\n", pkg.Header.SourceRpm())
	fmt.Printf("Summary:    %s\n", pkg.Header.Summary())
	fmt.Printf("Description:\n %s\n", pkg.Header.Description())

	return
}

func PrintPackagedFiles(file *os.File) (err error) {
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

func PrintPackageFileOf(file *os.File, filetype int32) (err error) {
	pkg, err := rpmlib.ReadPackageFile(file)
	if err != nil {
		return
	}

	filenames, err := pkg.Header.FileList()
	if err != nil {
		return
	}

	flags := pkg.Header.FileFlags()

	if len(flags) != len(filenames) {
		return fmt.Errorf("Included filenames size and flag size are different")
	}

	for i, filename := range filenames {
		if flags[i]&filetype != 0 {
			fmt.Println(filename)
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

type Option struct {
	ShowInfoMode       bool
	ShowFileMode       bool
	ShowConfigFileMode bool
	ShowDocFileMode    bool
	ShowChangelogMode  bool
	VerificationMode   bool
	CheckSignatureMode bool
}

func addOption(option *Option) {
	flag.BoolVar(&option.ShowInfoMode, "i", false, "Show package inforamtion.")
	flag.BoolVar(&option.ShowFileMode, "l", false, "Show files included package.")
	flag.BoolVar(&option.ShowConfigFileMode, "c", false, "Show config files included package.")
	flag.BoolVar(&option.ShowDocFileMode, "d", false, "Show doc files included package.")
	flag.BoolVar(&option.ShowChangelogMode, "changelog", false, "Show changelog.")
//	flag.BoolVar(&option.VerificationMode, "V", false,
//		"Verify file's size, checksum, permission, type, user and group.")
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
//		} else if option.VerificationMode {
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		file.Close()
	}

	os.Exit(0)
}
