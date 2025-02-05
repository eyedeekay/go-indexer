package indexer

import (
	"os"
	"path/filepath"
    "strings"
)

func (s *Scanner) Generate() error {
	if err := s.scan(); err != nil {
		return err
	}

	if err := s.parseTorrents(); err != nil {
		return err
	}

	return s.generateHTML()
}

// buildDirectoryTree creates a tree structure from scanned files
func (s *Scanner) buildDirectoryTree() *DirectoryNode {
	root := &DirectoryNode{
		Name:     filepath.Base(s.config.RootDir),
		Children: make(map[string]*DirectoryNode),
	}

	for _, file := range s.files {
		dir := filepath.Dir(file.RelPath)
		parts := strings.Split(dir, string(filepath.Separator))

		current := root
		for _, part := range parts {
			if part == "." {
				continue
			}
			if _, exists := current.Children[part]; !exists {
				current.Children[part] = &DirectoryNode{
					Name:     part,
					Children: make(map[string]*DirectoryNode),
				}
			}
			current = current.Children[part]
		}
		current.Files = append(current.Files, file)
	}

	return root
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
