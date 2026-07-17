APP := cldl
BUILD_DIR := build
BIN := $(BUILD_DIR)/$(APP)

GOFILES := $(shell find . -path './pkgbuild' -prune -o -name '*.go' -print)

PREFIX ?= /usr
DESTDIR ?=
BINDIR := $(DESTDIR)$(PREFIX)/bin
DATADIR := $(DESTDIR)$(PREFIX)/share/$(APP)

LDFLAGS := -ldflags "-w -s -X main.version=stripped -buildid= -extldflags=static"
EXTFLAGS := -buildvcs=false -a -installsuffix cgo -trimpath

.PHONY: build install uninstall clean

build: $(BIN)

$(BIN): $(GOFILES) go.mod go.sum
	@printf "\033[36m==> \033[0mCreating required directories...\n"
	@mkdir -p $(BUILD_DIR)

	@printf "\033[36m==> \033[0mBuilding $(APP)...\n"
	go build $(LDFLAGS) $(EXTFLAGS) -o $@ ./
	@printf "[\033[32m OK \033[0m] Build complete\n"

install: $(BIN)
	mkdir -p $(BINDIR)
	mkdir -p $(DATADIR)

	@printf "\033[36m==> \033[0mInstalling files...\n"
	install -m755 $(BIN) $(BINDIR)/$(APP)
	install -m644 default_config.toml $(DATADIR)/default_config.toml
	@printf "[\033[32m OK \033[0m] Installation complete\n"

uninstall:
	@printf "\033[36m==> \033[0mRemoving installed files...\n"
	rm -f $(PREFIX)/bin/$(APP)
	rm -rf $(DATADIR)
	@printf "[\033[32m OK \033[0m] Uninstall complete\n"

clean:
	@printf "\033[36m==> \033[0mRemoving build artifacts...\n"
	rm -rf $(BUILD_DIR)
	@printf "[\033[32m OK \033[0m] Clean complete\n"
