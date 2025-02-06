package backups

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func Compress(src, dst string) error {

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

	tarWalker := func(path string, fi os.FileInfo, err error) error {

		if fi.Mode().IsDir() {
			return nil
		}

		// this takes care of linked files that are problematic for tar
		if !fi.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if len(relPath) == 0 {
			return nil
		}

		rcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer rcFile.Close()

		if h, err := tar.FileInfoHeader(fi, relPath); err != nil {
			return err
		} else {
			h.Name = relPath
			if err = tw.WriteHeader(h); err != nil {
				return err
			}
		}

		if _, err := io.Copy(tw, rcFile); err != nil {
			return err
		}
		return nil
	}

	if err = filepath.Walk(src, tarWalker); err != nil {
		return err
	}

	return nil

}
