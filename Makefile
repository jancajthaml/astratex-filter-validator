.PHONY: all
all: vendor build

.ONESHELL:
.PHONY: windows
.PHONY: darwin

.PHONY: vendor
vendor:
	cd src && \
		go mod verify && \
		go mod tidy && \
		go mod vendor

.PHONY: build
build:
	$(MAKE) build-darwin
	$(MAKE) build-windows

.PHONY: build-%
build-%: %
	mkdir -p dist
	cd src && \
	\
	GOOS=$^ \
	GOARCH=amd64 \
	CGO_ENABLED=0 \
	GOFLAGS=$(shell [ "$^" = "windows" ] && echo "" || echo "-buildmode=pie") \
	\
	go build -a -o ../dist/$(shell [ "$^" = "windows" ] && echo "astratex-filter-validator.exe" || echo "astratex-filter-validator")

