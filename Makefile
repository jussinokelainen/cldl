APP := cldl
LDFLAGS := -ldflags "-w -s -X main.version=stripped -buildid= -extldflags=static"
EXTFLAGS := -buildvcs=false -a -installsuffix cgo -trimpath

build:
	@printf "\e[36m==> \e[0mBuilding $(APP)...\n"
	@printf "\e[36m==> \e[0mCreating required directories...\n"
	@mkdir -p bin

	@printf "\e[36m==> \e[0mCompiling binaries...\n"
	go build -o ./bin/$(APP) ./
	@printf "[\e[32m OK \e[0m] Debug Build complete\n"

release:
	@printf "\e[36m==> \e[0mBuilding $(APP)...\n"
	@printf "\e[36m==> \e[0mCreating required directories...\n"
	@mkdir -p bin

	@printf "\e[36m==> \e[0mCompiling binaries...\n"
	go build $(LDFLAGS) $(EXTFLAGS) -o ./bin/$(APP) ./
	@printf "[\e[32m OK \e[0m] Release Build complete\n"


# Install probably isn't the correct name of this, it is just there
# so i can easily copy the binary into my path on machines not on arch linux
install: release
	@printf "\e[32m==> \e[0mCopying required files and binaries...\n"
	cp ./bin/$(APP) ~/bin/$(APP)

	@printf "[\e[32m OK \e[0m] Installing complete\n"

clean:
	@printf "\e[32m==> \e[0mRemoving files and binaries...\n"
	rm -rf ./bin
	rm ~/bin/$(APP)
	@printf "[\e[32m OK \e[0m] Removing complete\n"
