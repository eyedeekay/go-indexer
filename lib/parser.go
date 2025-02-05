package indexer

import (
    "fmt"
    "sync"

    "github.com/xgfone/go-bt/metainfo"
)

type TorrentMeta struct {
    Name     string
    Size     int64
    Files    []TorrentFileInfo
    InfoHash string
}

type TorrentFileInfo struct {
    Path string
    Size int64
}

func (s *Scanner) parseTorrents() error {
    var wg sync.WaitGroup
    errs := make(chan error, len(s.files))

    for i := range s.files {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            if err := s.parseTorrent(&s.files[idx]); err != nil {
                errs <- fmt.Errorf("parsing %s: %w", s.files[idx].Path, err)
            }
        }(i)
    }

    wg.Wait()
    close(errs)

    // Collect errors
    var parseErrs []error
    for err := range errs {
        parseErrs = append(parseErrs, err)
    }

    if len(parseErrs) > 0 {
        return fmt.Errorf("parsing errors: %v", parseErrs)
    }
    return nil
}

func (s *Scanner) parseTorrent(tf *TorrentFile) error {
    mi, err := metainfo.LoadFromFile(tf.Path)
    if err != nil {
        return err
    }

    info, err := mi.Info()
    if err != nil {
        return err
    }

    meta := &TorrentMeta{
        Name:     info.Name,
        InfoHash: mi.InfoHash().String(),
    }

    if info.IsDir() {
        for _, file := range info.Files {
            meta.Files = append(meta.Files, TorrentFileInfo{
                Path: file.Path(info),
                Size: file.Length,
            })
            meta.Size += file.Length
        }
    } else {
        meta.Size = info.Length
        meta.Files = []TorrentFileInfo{{
            Path: info.Name,
            Size: info.Length,
        }}
    }

    tf.MetaInfo = meta
    return nil
}