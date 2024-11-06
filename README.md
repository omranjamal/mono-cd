# mono-cd

![GitHub Release Date](https://img.shields.io/github/release-date/omranjamal/mono-cd)
![GitHub Release](https://img.shields.io/github/v/release/omranjamal/mono-cd)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/omranjamal/mono-cd)


> The fastest way to `cd` into a JavaScript monorepo directory of interest

![DEMO](https://raw.githubusercontent.com/omranjamal/mono-cd/refs/heads/static/demo.gif)

## Features

1. `cd` into any of the folders in a JavaScript monorepo
2. Support for both interactive and non-interactive modes
3. Support for
   1. [pnpm workspaces](https://pnpm.io/workspaces)
   2. [npm workspaces](https://docs.npmjs.com/cli/v7/using-npm/workspaces/)
   3. [yarn workspaces](https://yarnpkg.com/features/workspaces)
   4. Custom directories via `.monocdrc.json`
4. Works inside Docker containers. (Tested with Alpine and Debian images)

## Usage

```bash
# interactive mode:
mcd

# non interactive `cd` if only one match is present.
mcd [search]
```

- `Up` / `Down` to select a target bookmark.
- Starting typing to filter list of bookmarks.

## Installation
Make sure you have `curl` installed on your system

```bash
curl -sL https://github.com/omranjamal/mono-cd/releases/latest/download/install.sh -o - | sh -
```

### Manual Installation

```bash
# Create installation directory
mkdir -p ~/.local/share/omranjamal/mono-cd

# Download the binary (check releases page for all available binaries)
curl -L -o ~/.local/share/omranjamal/mono-cd/mono-cd https://github.com/omranjamal/mono-cd/releases/latest/download/mono-cd_v1.1.0_amd64

# Add execution permissions
chmod +x ~/.local/share/omranjamal/mono-cd/mono-cd

# Add to shell
~/.local/share/omranjamal/mono-cd/mono-cd --install ~/.bashrc
```

### Setting Different Alias

You can either change the function name in your
`~/.bashrc` / `~/.zshrc` file from `mcd` to something
else.

OR, you could add the alias in Step 4 from above by passing
as the last argument.

```bash
~/.local/share/omranjamal/mono-cd/mono-cd --install ~/.bashrc monocd
```

`monocd` being the different alias that you want.

## Development

```bash
git clone git@github.com:omranjamal/mono-cd.git

cd ./mono-cd

go get
go run main.go
```

## License

MIT