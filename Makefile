BIN      := fzt
GOFLAGS  := -trimpath
LDFLAGS  := -s -w
# China-friendly proxy; override with `make GOPROXY=...` if needed.
GOPROXY  ?= https://goproxy.cn,direct
PREFIX   ?= /usr/local

.PHONY: all build install test clean fmt vet

all: build

build:
	GOPROXY=$(GOPROXY) go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN) .

install: build
	install -m 0755 $(BIN) $(DESTDIR)$(PREFIX)/bin/$(BIN)

test: build
	@mkdir -p demo/src/utils demo/docs demo/tests
	@touch demo/src/main.py demo/src/utils/helper.py demo/src/utils/config.py \
		demo/docs/README.md demo/tests/test_main.py
	python3 test_pty.py nav
	python3 test_pty.py search
	python3 test_pty.py escback
	@find demo -type f -delete
	@find demo -depth -type d -exec rmdir {} \;

# Cross-compiled release archives in dist/ (no cgo involved, so this is safe).
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
TARGETS  := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: release
release:
	@mkdir -p dist
	@for t in $(TARGETS); do \
		os=$${t%/*}; arch=$${t#*/}; \
		name=$(BIN)-$(VERSION)-$$os-$$arch; \
		ext=""; [ "$$os" = windows ] && ext=".exe"; \
		echo "==> $$name"; \
		mkdir -p dist/$$name; \
		GOOS=$$os GOARCH=$$arch go build $(GOFLAGS) -ldflags '$(LDFLAGS)' \
			-o dist/$$name/$(BIN)$$ext . || exit 1; \
		cp README.md key-bindings.bash key-bindings.zsh dist/$$name/; \
		if [ "$$os" = windows ]; then \
			(cd dist && zip -qr $$name.zip $$name); \
		else \
			(cd dist && tar czf $$name.tar.gz $$name); \
		fi; \
		rm -r dist/$$name; \
	done
	@cd dist && sha256sum *.tar.gz *.zip > checksums.txt 2>/dev/null; true
	@ls dist

fmt:
	gofmt -w *.go

vet:
	go vet ./...

clean:
	rm -f $(BIN)
	rm -rf dist
