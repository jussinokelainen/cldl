CLDL
====
Directory specific command line todo-lists.

Installation
------------
#### Arch Linux:
```bash
cd pkgbuild/arch
makepkg -si
```

#### macOS (homebrew):
```bash
brew tap jussinokelainen/cldl
brew trust jussinokelainen/cldl
brew install cldl
```

#### Other:
Run
```bash
make release
```
then copy the built binary from bin/ to your PATH.

Alternatively if ~/bin/ is in your PATH, run
```bash
make install
```
which builds the binary and copies it to ~/bin/
