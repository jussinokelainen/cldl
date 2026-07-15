APP := cldl
BUILD_DIR := build
LDFLAGS := -ldflags "-w -s -X main.version=stripped -buildid= -extldflags=static"
EXTFLAGS := -buildvcs=false -a -installsuffix cgo -trimpath

.PHONY: build release install clean uninstall

build:
	@printf "\033[36m==> \033[0mBuilding $(APP)...\n"
	@printf "\033[36m==> \033[0mCreating required directories...\n"
	@mkdir -p $(BUILD_DIR)

	@printf "\033[36m==> \033[0mCompiling binaries...\n"
	go build -o ./$(BUILD_DIR)/$(APP) ./
	@printf "[\033[32m OK \033[0m] Debug Build complete\n"

release:
	@printf "\033[36m==> \033[0mBuilding $(APP)...\n"
	@printf "\033[36m==> \033[0mCreating required directories...\n"
	@mkdir -p $(BUILD_DIR)

	@printf "\033[36m==> \033[0mCompiling binaries...\n"
	go build $(LDFLAGS) $(EXTFLAGS) -o ./$(BUILD_DIR)/$(APP) ./
	@printf "[\033[32m OK \033[0m] Release Build complete\n"


install: release
	@printf "\033[32m==> \033[0mCopying required files and binaries...\n"
	mkdir -p /usr/share/cldl
	cp ./$(BUILD_DIR)/$(APP) /usr/bin/$(APP)
	cp ./default_config.toml /usr/share/cldl/default_config.toml

	@printf "[\033[32m OK \033[0m] Installing complete\n"

uninstall:
	rm /usr/bin/$(APP)
	rm -rf /usr/share/cldl

clean:
	@printf "\033[32m==> \033[0mRemoving files and binaries...\n"
	rm -rf ./$(BUILD_DIR)
	@printf "[\033[32m OK \033[0m] Removing complete\n"
