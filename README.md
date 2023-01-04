# Go Snake Go
A snake game written in Go

![snake](https://user-images.githubusercontent.com/7248260/149492282-6588ead3-954d-42a4-9871-dc08cf833920.gif)

## Installation
### Get the binary
You can get the binary executable for your operating system from the [release page](https://github.com/liweiyi88/gosnakego/releases)

Windows: `gosnakego_windows_amd64.exe`  
MacOS: `gosnakego_darwin_amd64`  
Linux: `gosnakego_linux_amd64`

For MacOS, you can also run the following command to download the binary.
```
$ sudo wget https://github.com/liweiyi88/gosnakego/releases/download/v0.2.0/gosnakego_darwin_amd64 -O /usr/local/bin/gosnakego
$ sudo chmod +x /usr/local/bin/gosnakego
```

For Linux, you can also run the following command to download the binary.
```
$ sudo wget https://github.com/liweiyi88/gosnakego/releases/download/v0.2.0/gosnakego_linux_amd64 -O /usr/local/bin/gosnakego
$ sudo chmod +x /usr/local/bin/gosnakego
```

### Build yourself
```
$ git clone https://github.com/liweiyi88/gosnakego.git
$ make install
```

## How to play
Start the game with the command
```
$ gosnakego
```

or start the game with silent mode
```
$ gosnakego --silent
```

and use arrow keys to control the direction.
