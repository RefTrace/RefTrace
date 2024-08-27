# Binary name
BINARY_NAME=reft

# Go build command
GOBUILD=go build

# Build flags
BUILD_FLAGS=-ldflags="-s -w"

# Targets
.PHONY: all build clean linux mac

all: linux mac

build: all

linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-x86_64

mac:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-mac-m1

clean:
	rm -f $(BINARY_NAME)-linux-x86_64 $(BINARY_NAME)-mac-m1