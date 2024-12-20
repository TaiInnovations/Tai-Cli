#!/bin/bash
mkdir -p build

# Darwin (macOS) builds
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/free-gemini_darwin_amd64
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/free-gemini_darwin_arm64

# Windows builds
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/free-gemini_windows_amd64.exe
GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/free-gemini_windows_arm64.exe
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/free-gemini_windows_386.exe

# Linux builds
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/free-gemini_linux_amd64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/free-gemini_linux_arm64
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/free-gemini_linux_386

# Set executable permissions
chmod +x build/free-gemini_darwin_amd64
chmod +x build/free-gemini_darwin_arm64
chmod +x build/free-gemini_windows_amd64.exe
chmod +x build/free-gemini_windows_arm64.exe
chmod +x build/free-gemini_windows_386.exe
chmod +x build/free-gemini_linux_amd64
chmod +x build/free-gemini_linux_arm64
chmod +x build/free-gemini_linux_386
