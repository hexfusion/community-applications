BUILD_DIR = build

BUILD_SCRIPT = build.go

build: clean
	@echo "Building the static site..."
	go run $(BUILD_SCRIPT)

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

help:
	@echo "Usage:"
	@echo "  make build  - Build the static site"
	@echo "  make clean  - Remove the build directory"
	@echo "  make help   - Show this help message"