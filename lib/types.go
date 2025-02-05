// lib/types.go
package indexer

import (
	"sync"
)

// Config holds scanner configuration
type Config struct {
	RootDir string
	Output  string
}

// Scanner handles torrent file scanning and HTML generation
type Scanner struct {
	config Config
	files  []TorrentFile
	mu     sync.Mutex
}

// DirectoryNode represents a node in the directory tree
type DirectoryNode struct {
	Name     string
	Path     string
	Files    []TorrentFile
	Children map[string]*DirectoryNode
}

// TorrentFile represents a torrent file with its metadata
type TorrentFile struct {
	Path     string
	RelPath  string
	MetaInfo *TorrentMeta
}

// TorrentMeta holds parsed torrent metadata
type TorrentMeta struct {
	Name     string
	Size     int64
	Files    []TorrentFileInfo
	InfoHash string
}

// TorrentFileInfo holds information about files within a torrent
type TorrentFileInfo struct {
	Path string
	Size int64
}

// NewScanner creates a new Scanner instance
func NewScanner(config Config) *Scanner {
	return &Scanner{
		config: config,
		files:  make([]TorrentFile, 0),
	}
}
