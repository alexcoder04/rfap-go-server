
LINUX_OUT = ../build/rfap-go-server_linux_amd64
RASPBERRY_OUT = ../build/rfap-go-server_linux_arm
WINDOWS_OUT = ../build/rfap-go-server_windows_amd64.exe

BUILD_TIMESTAMP = $(shell date --utc "+%Y.%m.%d-%H:%M:%S")
GIT_COMMIT = $(shell git rev-list -1 HEAD)
VERSION = $(shell git describe --tags --abbrev=0)
BUILD_OS = $(shell uname -orm | tr "[:blank:]" "_")

LDFLAGS = "-X main.BUILD_TIMESTAMP=$(BUILD_TIMESTAMP) -X main.GIT_COMMIT=$(GIT_COMMIT) -X main.SERVER_VERSION=$(VERSION) -X main.BUILD_OS=$(BUILD_OS)"

linux:
	GOOS=linux GOARCH=amd64\
		 go build -ldflags $(LDFLAGS)\
		 -o $(LINUX_OUT)

windows:
	GOOS=windows GOARCH=amd64\
		 go build -ldflags $(LDFLAGS)\
		 -o $(WINDOWS_OUT)

run:
	RFAP_MODE=testing RFAP_LOG_LEVEL=trace\
		go run -ldflags $(LDFLAGS) .

install:
	go install -ldflags $(LDFLAGS) .

clean:
	rm -f $(LINUX_OUT) $(RASPBERRY_OUT) $(WINDOWS_OUT)

