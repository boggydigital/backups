package konpo

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
)

func Compress(src, dst string) error {

	if src == "" || dst == "" {
		return errors.New("compressing requires src and dst dirs")
	}

	exportedPath := filepath.Join(dst, Filename())

	if _, err := os.Stat(exportedPath); os.IsExist(err) {
		return err
	}

	file, err := os.Create(exportedPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	fs := os.DirFS(src)

	if err := tw.AddFS(fs); err != nil {
		return err
	}

	return nil
}
