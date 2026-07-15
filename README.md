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
```bash
sudo make install
```

#### Uninstall:
Uninstalling can be done through your package manager
if it was used for installing, otherwise
```bash
sudo make uninstall
```
