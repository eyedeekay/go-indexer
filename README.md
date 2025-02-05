# go-indexer

A command-line tool written in Go that recursively scans directories for torrent files and generates a searchable HTML index.

## Features

- Scans directories recursively for .torrent files
- Extracts metadata including file sizes and info hashes
- Generates a searchable HTML index with:
  - Collapsible directory structure
  - File listings for each torrent
  - Human-readable file sizes
  - Real-time search functionality

## Installation

```bash
go install github.com/eyedeekay/go-indexer@latest
```

## Usage

```bash
# Basic usage - scans current directory
go-indexer

# Specify custom directory and output file
go-indexer -dir /path/to/torrents -output custom.html
```

### Flags

- `-dir` Root directory to scan (default: current directory)
- `-output` Output HTML file path (default: index.html)

## Build Requirements

- Go 1.20 or later
- github.com/xgfone/go-bt v0.6.1

## License

MIT License - See LICENSE file