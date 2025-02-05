// main.go
package main

import (
    "flag"
    "log"
    "path/filepath"

    "github.com/eyedeekay/go-indexer"
)

func main() {
    rootDir := flag.String("dir", ".", "Root directory containing torrent files")
    output := flag.String("output", "index.html", "Output HTML file path")
    flag.Parse()

    absRoot, err := filepath.Abs(*rootDir)
    if err != nil {
        log.Fatalf("Invalid root directory: %v", err)
    }

    scanner := indexer.NewScanner(indexer.Config{
        RootDir: absRoot,
        Output:  *output,
    })

    if err := scanner.Generate(); err != nil {
        log.Fatalf("Generation failed: %v", err)
    }
}