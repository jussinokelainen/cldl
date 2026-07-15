CLDL
====
Directory specific command line todo-lists.

Installation
------------
#### Arch Linux:
```bash
git clone git@github.com:jussinokelainen/cldl.git
cd cldl/pkgbuild/arch
makepkg -si
```

#### Void Linux:
```bash
git clone git@github.com:jussinokelainen/cldl.git
cd cldl/pkgbuild/void
./install
```

#### Debian:
```bash
git clone git@github.com:jussinokelainen/cldl.git
cd cldl/pkgbuild/debian
./install
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
git clone git@github.com:jussinokelainen/cldl.git
cd cldl
make release
```
then copy the built binary from build/ to your PATH.

Alternatively if ~/bin/ is in your PATH, run
```bash
make install
```
which builds the binary and copies it to ~/bin/
