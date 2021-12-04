
LINUX_OUT = rfap-go-server_linux_amd64
RASPBERRY_OUT = rfap-go-server_linux_arm
WINDOWS_OUT = rfap-go-server_windows_amd64.exe

linux:
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_OUT)

raspberry:
	GOOS=linux GOARCH=arm go build -o $(RASPBERRY_OUT)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_OUT)

run:
	RFAP_MODE=testing RFAP_LOG_LEVEL=trace go run .

clean:
	rm -f $(LINUX_OUT) $(RASPBERRY_OUT) $(WINDOWS_OUT)

