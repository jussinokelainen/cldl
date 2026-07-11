APP := todo
LDFLAGS := -ldflags "-w -s -X main.version=stripped -buildid= -extldflags=static"
EXTFLAGS := -buildvcs=false -a -installsuffix cgo -trimpath

# This makefile and program will probably not work as intended
# if directories and personal project stuff etc aren't the same

# NOTE: Keep this as the first
build:
	@printf "\e[36m==> \e[0mBuilding $(APP)...\n"
	@printf "\e[36m==> \e[0mCreating required directories...\n"
	@mkdir -p bin

	@printf "\e[36m==> \e[0mCompiling binaries...\n"
	go build -o ./bin/$(APP) ./

	@printf "[\e[32m OK \e[0m] Build complete\n"

# Remove all files created by the program
clean:
	@printf "\e[32m==> \e[0mRemoving files and binaries...\n"
	rm -rf ./bin
	rm ~/bin/$(APP)
	@printf "[\e[32m OK \e[0m] Removing complete\n"

install: build
	@printf "\e[32m==> \e[0mCopying required files and binaries...\n"
	go build $(LDFLAGS) $(EXTFLAGS) -o ./bin/$(APP)-release ./
	cp ./bin/$(APP)-release ~/bin/$(APP)

	@printf "[\e[32m OK \e[0m] Installing complete\n"
