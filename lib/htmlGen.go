// lib/htmlGen.go
package indexer

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// TemplateData holds the data for HTML template
type TemplateData struct {
	DirectoryTree template.HTML
}

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

// generateHTML creates the final HTML file
func (s *Scanner) generateHTML() error {
	tree := s.buildDirectoryTree()
	htmlContent := s.renderDirectoryTree(tree)

	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(s.config.Output)
	if err != nil {
		return err
	}
	defer f.Close()

	data := TemplateData{
		DirectoryTree: template.HTML(htmlContent),
	}

	return tmpl.Execute(f, data)
}

// renderDirectoryTree converts the directory tree to HTML
func (s *Scanner) renderDirectoryTree(node *DirectoryNode) string {
	var sb strings.Builder

	id := fmt.Sprintf("dir_%p", node)

	if node.Name != "" {
		sb.WriteString("<div class=\"directory\">")
		sb.WriteString(fmt.Sprintf("<input type=\"checkbox\" id=\"%s\" checked>", id))
		sb.WriteString(fmt.Sprintf("<label for=\"%s\">📁 %s</label>", id, template.HTMLEscapeString(node.Name)))
		sb.WriteString("<div class=\"contents\">")
	}

	for _, file := range node.Files {
		sb.WriteString("<div class=\"torrent-info\">")
		sb.WriteString(fmt.Sprintf("📄 <strong>%s</strong><br>", template.HTMLEscapeString(filepath.Base(file.RelPath))))
		if file.MetaInfo != nil {
			sb.WriteString(fmt.Sprintf("Size: %s<br>", formatSize(file.MetaInfo.Size)))
			sb.WriteString(fmt.Sprintf("InfoHash: %s<br>", template.HTMLEscapeString(file.MetaInfo.InfoHash)))
			if len(file.MetaInfo.Files) > 0 {
				sb.WriteString("<details><summary>Files:</summary><ul>")
				for _, f := range file.MetaInfo.Files {
					sb.WriteString(fmt.Sprintf("<li>%s (%s)</li>",
						template.HTMLEscapeString(f.Path),
						formatSize(f.Size)))
				}
				sb.WriteString("</ul></details>")
			}
		}
		sb.WriteString("</div>\n")
	}

	for _, child := range node.Children {
		sb.WriteString(s.renderDirectoryTree(child))
	}

	if node.Name != "" {
		sb.WriteString("</div></div>")
	}

	return sb.String()
}

// formatSize converts bytes to human-readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
