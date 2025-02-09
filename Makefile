# Define build directory
BUILD_DIR := build

# Define output names for different platforms
DARWIN_AMD64 := $(BUILD_DIR)/Tai_darwin_amd64
DARWIN_ARM64 := $(BUILD_DIR)/Tai_darwin_arm64
WINDOWS_AMD64 := $(BUILD_DIR)/Tai_windows_amd64.exe
WINDOWS_ARM64 := $(BUILD_DIR)/Tai_windows_arm64.exe
WINDOWS_386 := $(BUILD_DIR)/Tai_windows_386.exe
LINUX_AMD64 := $(BUILD_DIR)/Tai_linux_amd64
LINUX_ARM64 := $(BUILD_DIR)/Tai_linux_arm64
LINUX_386 := $(BUILD_DIR)/Tai_linux_386

# Default target, build all platforms
.PHONY: all
all: $(DARWIN_AMD64) $(DARWIN_ARM64) $(WINDOWS_AMD64) $(WINDOWS_ARM64) $(WINDOWS_386) $(LINUX_AMD64) $(LINUX_ARM64) $(LINUX_386)

# Darwin builds
$(DARWIN_AMD64):
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $@

$(DARWIN_ARM64):
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $@

# Windows builds
$(WINDOWS_AMD64):
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $@

$(WINDOWS_ARM64):
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $@

$(WINDOWS_386):
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $@

# Linux builds
$(LINUX_AMD64):
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $@

$(LINUX_ARM64):
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $@

$(LINUX_386):
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $@

# Set executable permissions
.PHONY: chmod
chmod: $(DARWIN_AMD64) $(DARWIN_ARM64) $(WINDOWS_AMD64) $(WINDOWS_ARM64) $(WINDOWS_386) $(LINUX_AMD64) $(LINUX_ARM64) $(LINUX_386)
	chmod +x $(DARWIN_AMD64)
	chmod +x $(DARWIN_ARM64)
	chmod +x $(WINDOWS_AMD64)
	chmod +x $(WINDOWS_ARM64)
	chmod +x $(WINDOWS_386)
	chmod +x $(LINUX_AMD64)
	chmod +x $(LINUX_ARM64)
	chmod +x $(LINUX_386)

# Clean up the build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Multi-platform build
.PHONY: build
build:
	$(MAKE) darwin_amd64 darwin_arm64 windows_amd64 windows_arm64 linux_amd64 linux_arm64

# Targets for specific platforms (you can run these individually)
.PHONY: darwin_amd64
darwin_amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $(DARWIN_AMD64)

.PHONY: darwin_arm64
darwin_arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $(DARWIN_ARM64)

.PHONY: windows_amd64
windows_amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $(WINDOWS_AMD64)

.PHONY: windows_arm64
windows_arm64:
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $(WINDOWS_ARM64)

.PHONY: linux_amd64
linux_amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $(LINUX_AMD64)

.PHONY: linux_arm64
linux_arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $(LINUX_ARM64)
