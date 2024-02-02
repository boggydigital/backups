package hogo

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"github.com/boggydigital/nod"
	"io"
	"os"
	"path/filepath"
)

func Compress(src, dst string, tpw nod.TotalProgressWriter) error {

	if src == "" || dst == "" {
		return errors.New("compressing requires src and dst dirs")
	}

	root, _ := filepath.Split(src)

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

	files := make([]string, 0)

	if err := filepath.Walk(src, func(f string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		files = append(files, f)
		return nil
	}); err != nil {
		return err
	}

	if tpw != nil {
		tpw.TotalInt(len(files))
	}

	for _, f := range files {

		fi, err := os.Stat(f)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, f)
		if err != nil {
			return err
		}

		rp, err := filepath.Rel(root, f)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(rp)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		of, err := os.Open(f)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, of); err != nil {
			return err
		}

		if tpw != nil {
			tpw.Increment()
		}
	}

	return nil
}
