PROG=bin/redisproxy


SRCS=./cmd/proxy


# git commit hash
COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")

# Compile the date
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')

# compile
CFLAGS = -ldflags "-s -w -X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\""

GOPROXY=https://goproxy.cn,direct

all:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	GOPROXY=$(GOPROXY) go build $(CFLAGS) -o $(PROG) $(SRCS)

# Compiling the RACE version
race:
	if [ ! -d "./bin/" ]; then \
    	mkdir bin; \
    	fi
	GOPROXY=$(GOPROXY) go build $(CFLAGS) -race -o $(PROG) $(SRCS)

# release version
RELEASE_DATE = $(shell date '+%Y%m%d%H%M%S')
RELEASE_VERSION = $(shell git rev-parse --short HEAD || echo "GitNotFound")
RELEASE_DIR=release_bin
RELEASE_BIN_NAME=RedisProxy
release:
	if [ ! -d "./$(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION)" ]; then \
	mkdir ./$(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION); \
	fi
	GOPROXY=$(GOPROXY) go build $(CFLAGS) -o $(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION)/$(RELEASE_BIN_NAME) $(SRCS)

clean:
	rm -rf ./bin

run:
	GOPROXY=$(GOPROXY) go run --race cmd/proxy/main.go -c config/config.yaml