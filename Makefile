# Binary and image names
BINARY      := myapp
IMAGE       := myapp:0.0.1
DOCKERFILE  := Dockerfile

# Default Go build flags
GOFLAGS     := -ldflags="-s -w" -trimpath
GOOS        := linux
GOARCH      := amd64

.PHONY: all build clean docker docker-run

all: build

# Build a static Linux binary
build:
	@echo ">> Building Go binary for $(GOOS)/$(GOARCH)"
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -o $(BINARY) .

# Remove build artifacts
clean:
	@echo ">> Cleaning build artifacts"
	rm -f $(BINARY)

# Build Docker image
docker: build
	@echo ">> Building Docker image: $(IMAGE)"
	docker build -t $(IMAGE) -f $(DOCKERFILE) .

# Run container locally (ports mapped)
docker-run:
	@echo ">> Running container: $(IMAGE)"
	docker run --rm -p 8080:8080 $(IMAGE)

# output the zarf package create command

key:
	@echo "zarf tools gen-key"

create:
	@echo "zarf package create --signing-key cosign.key --signing-key-pass test"

sbom:
	@echo "zarf package inspect sbom zarf-package-cnad-amd64-0.0.1.tar.zst --output sbom"

deploy:
	@echo "zarf package deploy zarf-package-cnad-amd64-0.0.1.tar.zst --key cosign.pub"

unarchive:
	@echo "zarf tools archiver decompress zarf-package-cnad-amd64-0.0.1.tar.zst archive --unarchive-all"