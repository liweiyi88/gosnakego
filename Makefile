.PHONY: install release

ARTIFACTS_DIR=artifacts/${VERSION}
GITHUB_USERNAME=liweiyi88

install:
	go install

release:
	GOOS=windows GOARCH=amd64 go build -o $(ARTIFACTS_DIR)gosnakego_windows_amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o $(ARTIFACTS_DIR)gosnakego_darwin_amd64
	GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)gosnakego_linux_amd64
