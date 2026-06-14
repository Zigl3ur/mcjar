# McJar

A CLI tool to download minecraft server jars, mods, plugins, modpacks, and datapacks.

## Installation

### Binary

  Grab it from the [releases](https://github.com/Zigl3ur/mcjar/releases/latest)

### Docker

```bash
docker pull ghcr.io/zigl3ur/mcjar:latest
docker run --rm -it ghcr.io/zigl3ur/mcjar:latest <args>
```

### Build from source
```bash
git clone https://github.com/Zigl3ur/mcjar.git
cd mcjar
go mod tidy
go build
```

## Commands

- `jar` - Download a minecraft server jar

  - `list` - List available versions / builds for a specified server type

- `addons` - Download mods / plugins / modpacks / datapacks from modrinth
  - `search` - Search for mods / plugins / modpacks / datapacks
  - `info` - Get info about a mod / plugin / modpack / datapack
  - `get` - Download a mod / plugin / modpack / datapack
