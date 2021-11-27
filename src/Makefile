
LINUX_OUT = rfap-go-server_linux_amd64
RASPBERRY_OUT = rfap-go-server_linux_arm64
WINDOWS_OUT = rfap-go-server_windows_amd64.exe

linux:
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_OUT)

raspberry:
	GOOS=linux GOARCH=arm64 go build -o $(RASPBERRY_OUT)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_OUT)

clean:
	rm -f $(LINUX_OUT) $(RASPBERRY_OUT) $(WINDOWS_OUT)

