CLDL
====
Directory specific command line todo-lists.

Configuring
-----------
Configuring cldl can be done through a config.toml file, which should
be located at $HOME/.config/cldl/config.toml

The file can be either created manually, or by running the command below
to generate a config file containing the default settings.
```bash
cldl generate-configs
```

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

Uninstalling
------------
If you wish to remove all saved lists (that are listed in 'cldl ls -a'), run
```bash
cldl delete-lists
```
before uninstalling the package.

Uninstalling can be done through your package manager, if it was
used for installing. Otherwise:
```bash
sudo make uninstall
```
Or manually:
```bash
sudo rm /usr/bin/cldl
sudo rm -rf /usr/share/cldl
```
