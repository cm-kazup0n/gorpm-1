.PHONY: all gorpm gorpm2cpio
all:
	make gorpm2cpio
	make gorpm
gorpm2cpio:
	go build -ldflags="-s -w" -o ./build/gorpm2cpio ./gorpm2cpio/gorpm2cpio.go
gorpm:
	go build -ldflags="-s -w" -o ./build/gorpm ./gorpm/gorpm.go