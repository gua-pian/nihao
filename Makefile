# CUR_TAG is the last git tag plus the delta from the current commit to the tag
# e.g. v1.5.5-<nr of commits since>-g<current git sha>
CUR_TAG = $(shell git describe)

# LAST_TAG is the last git tag
# e.g. v1.5.5
LAST_TAG = $(shell git describe --abbrev=0)

# VERSION is the last git tag without the 'v'
# e.g. 1.5.5
VERSION = $(shell git describe --abbrev=0 | cut -c 2-)

# GO runs the go binary with garbage collection disabled for faster builds.
# Do not specify a full path for go since travis will fail.
GO = go

# GOFLAGS is the flags for the go compiler. Currently, only the version number is
# passed to the linker via the -ldflags.
GOFLAGS = -ldflags "-X main.version=$(CUR_TAG)"

# GOVERSION is the current go version, e.g. go1.9.2
GOVERSION = $(shell go version | awk '{print $$3;}')


# all is the default target
all: test

# help prints a help screen
help:
	@echo "build     - go build (for current platform)"
	@echo "install   - go install (for current platform)"
	@echo "gofmt     - go fmt"
	@echo "linux     - go build linux/amd64 (cross compile)"
	@echo "windows   - go build windows/386 (cross compile)"
	@echo "macos     - go build darwin/amd64 (cross compile)"
	@echo "clean     - remove temp files"

# build compiles fabio and the test dependencies
build:  gofmt
	$(GO) build -tags netgo $(GOFLAGS)

# gofmt runs gofmt on the code
gofmt:
	gofmt -s -w `find . -type f -name '*.go' | grep -v vendor`

# linux builds a linux binary
linux:
	GOOS=linux GOARCH=amd64 $(GO) build -tags netgo $(GOFLAGS)

# windows builds a exe file
windows:
	GOOS=windows GOARCH=386 $(GO) build -tags netgo $(GOFLAGS)

macos:
	GOOS=darwin GOARCH=amd64 $(GO) build -tags netgo $(GOFLAGS)


# install runs go install
install:
	$(GO) install $(GOFLAGS)


# clean removes intermediate files
clean:
	$(GO) clean
	rm -rf pkg dist awesomeProject awesomeProject.exe
	find . -name '*.test' -delete

