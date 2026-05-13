# MockResponder Go Plugin — Build Targets
#
# Pre-plugin that returns OAS example responses for mock requests.
# Enables portal "Try It" without credentials or live backend.

TYK_VERSION ?= v5.11.0
PLUGIN_NAME := mockResponder
SO_FILE     := $(PLUGIN_NAME).so

.PHONY: build build-arm64 test clean

# Compile the Go plugin (AMD64) using tyk-plugin-compiler.
build:
	docker run --rm \
		--platform linux/amd64 \
		-v $(CURDIR):/plugin-source \
		tykio/tyk-plugin-compiler:$(TYK_VERSION) $(SO_FILE)
	@echo "Built $(SO_FILE) for Tyk $(TYK_VERSION) (amd64)"

# Compile the Go plugin (ARM64/Graviton) using tyk-plugin-compiler.
build-arm64:
	docker run --rm \
		--platform linux/arm64 \
		-v $(CURDIR):/plugin-source \
		tykio/tyk-plugin-compiler:$(TYK_VERSION) $(SO_FILE)
	@echo "Built $(SO_FILE) for Tyk $(TYK_VERSION) (arm64)"

# Run unit tests inside the plugin compiler container.
test:
	docker run --rm \
		--platform linux/amd64 \
		-v $(CURDIR):/plugin-source \
		-w /plugin-source \
		--entrypoint "" \
		tykio/tyk-plugin-compiler:$(TYK_VERSION) \
		bash -c "go test -v -count=1 ./..."

# Remove build artifacts.
clean:
	rm -f *.so
