# Mcli

A CLI tool to download minecraft server jars, mods, plugins, modpacks, datapacks and create backups of server folders.

## Installation

Download the latest release [here](https://github.com/Zigl3ur/mcli/releases).

Or clone the repo and build it:

```bash
git clone https://github.com/Zigl3ur/mcli.git
cd mcli
go build
```

## Commands

- `jar` - Download a minecraft server jar

  - `list` - List available versions / builds for a specified server type

- `addons` - Download mods / plugins / modpacks / datapacks from modrinth

  - `search` - Search for mods / plugins / modpacks / datapacks
  - `info` - Get info about a mod / plugin / modpack / datapack
  - `get` - Download a mod / plugin / modpack / datapack

- `backup` - Create a backup of a specified server folder

## TODOS

### Subcommands

- [x] Implement `list`, to list versions / build for specified server type
- [x] Implement `jar`, to download server jar from specified servertype / version / build
- [x] implement `addons`, to download mods / plugins / modpacks / datapacks
- [x] Implement `backup`, to create a backup of the specified server folder

### Misc

- [x] Function to sort an array of minecraft versions
- [x] Display download speed
- [ ] Add Tests for handlers and commands
- [x] change default output if not specified to smthng like `neoforge-1.21.8-234.jar` (serverType-version-build.jar)
- [x] for mods/plugins/modpacks/datapacks search display filter used (if any)
- [ ] Better errors messages

```

```
