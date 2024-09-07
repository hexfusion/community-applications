BUILD_DIR = build

BUILD_SITE_SCRIPT = cmd/build-site/build-site.go
GENERATE_TEMPLATES_SCRIPT = cmd/generate-templates/generate-templates.go

build-site: clean
	@echo "Building the static site..."
	go run $(BUILD_SITE_SCRIPT)

generate:
	@echo "Generating templates..."
	go run $(GENERATE_TEMPLATES_SCRIPT)

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

help:
	@echo "Usage:"
	@echo "  make build-site  - Build the static site"
	@echo "  make clean  - Remove the build directory"
	@echo "  make help   - Show this help message"