gorpm
===

RPM commmand line utility implemented in Go language

## Description
RPM commmand line utility implemented in Go language.
This 'gorpm' and other commands aim to get RPM package's 
information at such as Windows, debian,which does not have yum/rpm command. 
So this command is not for package management on RHEL or CentOS.

This is my first 'decent' Go language project for practicing.
If you found any mistakes or bugs in command, please tell me.

## Installation

```
$ git clone　git@github.com:cm-kazup0n/gorpm-1.git
$ cd gorpm-1
$ make
$ ls build/
gorpm  gorpm2cpio
```

## Usages

### gorpm

* Show package information

``` 
$ gorpm -i <RPM Package>
```

```
$ gorpm -i rpm-4.8.0-55.el6.x86_64.rpm
Name:       rpm
Version:    4.8.0
Release:    55.el6
Group:      System Environment/Base
Size:       2034245
Licence:    GPLv2+
BuildDate:  2016-05-11 08:49:46 +0900 JST
Source RPM: rpm-4.8.0-55.el6.src.rpm
Summary:    The RPM package management system
Description:
 The RPM Package Manager (RPM) is a powerful command line driven
 package management system capable of installing, uninstalling,
 verifying, querying, and updating software packages. Each software
 package consists of an archive of files along with information about
 the package like its version, a description, etc.
```

* Show files in RPM package

``` 
$ gorpm -l <RPM Package>
```

```
$ gorpm -i rpm-4.8.0-55.el6.x86_64.rpm
/bin/rpm
/etc/rpm
/usr/bin/rpm2cpio
/usr/bin/rpmdb
/usr/bin/rpmquery
/usr/bin/rpmsign
/usr/bin/rpmverify
/usr/lib/rpm
/usr/lib/rpm/macros
/usr/lib/rpm/platform
/usr/lib/rpm/platform/amd64-linux
/usr/lib/rpm/platform/amd64-linux/macros
~ ~ ~
```

* Verify RPM Package

``` 
$ gorpm -V <RPM Package>
```

This option is not full compatibility for `rpm -qV`
Currently, verify file's size, mode, mtime and checksum.


### gorpm2cpio
Reimplementaion of rpm2cpio command.

``` 
$ gorpm2cpio <RPM Package>
```

CPIO archive data is dumped to stdout.
So if you want take out files from archive, in the same way rpm2cpio,
hit following command.

```
$ gorpm2cpio | cpio -id
```

## FAQ
1. Why xx option has not been implemented ? When will you implement it ?
 Sometime when I need it. Or sometime when others give me an early Xmas present.


## Author
Necomeshi
