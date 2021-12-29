
LINUX_OUT = build/rfap-go-server_linux_amd64
WINDOWS_OUT = build/rfap-go-server_windows_amd64.exe
MAC_INTEL_OUT = build/rfap-go-server_darwin_amd64
RASPBERRY_OUT = build/rfap-go-server_linux_arm

BUILD_TIMESTAMP = $(shell date --utc "+%Y.%m.%d-%H:%M:%S")
GIT_COMMIT = $(shell git rev-list -1 HEAD)
VERSION = $(shell git describe --tags --abbrev=0)
BUILD_OS = $(shell uname -orm | tr "[:blank:]" "_")

LDFLAGS = "-X github.com/alexcoder04/rfap-go-server/settings.BUILD_TIMESTAMP=$(BUILD_TIMESTAMP) -X github.com/alexcoder04/rfap-go-server/settings.GIT_COMMIT=$(GIT_COMMIT) -X github.com/alexcoder04/rfap-go-server/settings.SERVER_VERSION=$(VERSION) -X github.com/alexcoder04/rfap-go-server/settings.BUILD_OS=$(BUILD_OS)"

linux:
	GOOS=linux GOARCH=amd64\
		 go build -ldflags $(LDFLAGS)\
		 -o $(LINUX_OUT) .

raspberry:
	GOOS=linux GOARCH=arm\
		 go build -ldflags $(LDFLAGS)\
		 -o $(RASPBERRY_OUT) .

windows:
	GOOS=windows GOARCH=amd64\
		 go build -ldflags $(LDFLAGS)\
		 -o $(WINDOWS_OUT) .

mac-intel:
	GOOS=darwin GOARCH=amd64\
		 go build -ldflags $(LDFLAGS)\
		 -o $(MAC_INTEL_OUT) .

run:
	RFAP_LOG_FILE=[stdout] RFAP_LOG_FORMAT=color RFAP_LOG_LEVEL=trace\
		go run -ldflags $(LDFLAGS) .

run-quiet:
	RFAP_LOG_FILE=[stdout] RFAP_LOG_LEVEL=info\
		go run -ldflags $(LDFLAGS) .

install:
	go install -ldflags $(LDFLAGS) .

clean:
	rm -f $(LINUX_OUT) $(WINDOWS_OUT) $(MAC_INTEL_OUT) $(RASPBERRY_OUT)

