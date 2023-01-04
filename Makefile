.PHONY: build release

ARTIFACTS_DIR=artifacts/${VERSION}
GITHUB_USERNAME=liweiyi88

build:
	go install
	go build .

release:
	GOOS=windows go build -o $(ARTIFACTS_DIR)gosnakego_windows_amd64.exe
	GOOS=darwin go build -o $(ARTIFACTS_DIR)gosnakego_darwin_amd64
