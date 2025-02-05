package indexer

import (
    "os"
    "path/filepath"
    "sync"
)

type Config struct {
    RootDir string
    Output  string
}

type Scanner struct {
    config Config
    files  []TorrentFile
    mu     sync.Mutex
}

type TorrentFile struct {
    Path     string
    RelPath  string
    MetaInfo *TorrentMeta
}

func NewScanner(config Config) *Scanner {
    return &Scanner{
        config: config,
        files:  make([]TorrentFile, 0),
    }
}

func (s *Scanner) Generate() error {
    if err := s.scan(); err != nil {
        return err
    }

    if err := s.parseTorrents(); err != nil {
        return err
    }

    return s.generateHTML()
}

func (s *Scanner) scan() error {
    return filepath.Walk(s.config.RootDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        if filepath.Ext(path) != ".torrent" {
            return nil
        }

        relPath, err := filepath.Rel(s.config.RootDir, path)
        if err != nil {
            return err
        }

        s.mu.Lock()
        s.files = append(s.files, TorrentFile{
            Path:    path,
            RelPath: relPath,
        })
        s.mu.Unlock()

        return nil
    })
}