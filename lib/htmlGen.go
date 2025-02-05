package indexer

import (
    "html/template"
    "os"
    "path/filepath"
    "strings"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Torrent Directory</title>
    <style>
        .directory { margin-left: 20px; }
        .directory label { cursor: pointer; }
        .directory input[type="checkbox"] { display: none; }
        .directory input:not(:checked) ~ .contents { display: none; }
        .torrent-info { margin: 5px 0; padding: 5px; border: 1px solid #ddd; }
        .search { margin: 20px 0; }
        .hidden { display: none; }
    </style>
    <script>
        function search() {
            const query = document.getElementById('search').value.toLowerCase();
            document.querySelectorAll('.torrent-info').forEach(el => {
                const text = el.textContent.toLowerCase();
                el.classList.toggle('hidden', !text.includes(query));
            });
        }
    </script>
</head>
<body>
    <div class="search">
        <input type="text" id="search" onkeyup="search()" placeholder="Search torrents...">
    </div>
    <div class="directory">
        {{.DirectoryTree}}
    </div>
</body>
</html>`

func (s *Scanner) generateHTML() error {
    tree := s.buildDirectoryTree()
    
    tmpl, err := template.New("index").Parse(htmlTemplate)
    if err != nil {
        return err
    }

    f, err := os.Create(s.config.Output)
    if err != nil {
        return err
    }
    defer f.Close()

    return tmpl.Execute(f, tree)
}

type DirectoryNode struct {
    Name     string
    Path     string
    Files    []TorrentFile
    Children map[string]*DirectoryNode
}

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