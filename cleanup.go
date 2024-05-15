package backups

import (
	"github.com/boggydigital/nod"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	daysToPreserveFiles = 30
)

func Cleanup(dir string, delete bool, tpw nod.TotalProgressWriter) error {

	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	filenames, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	earliest := time.Now().Add(-daysToPreserveFiles * 24 * time.Hour)
	oldFiles := make([]string, 0)

	for _, fn := range filenames {

		fnse := fn
		for filepath.Ext(fnse) != "" {
			fnse = strings.TrimSuffix(fnse, filepath.Ext(fnse))
		}
		ft, err := time.Parse(nod.TimeFormat, fnse)
		if err != nil {
			nod.Log(err.Error())
			continue
		}

		if ft.After(earliest) {
			continue
		}

		oldFiles = append(oldFiles, fn)
	}

	if len(oldFiles) > 0 && delete {

		// never delete all backups, leave the latest file as the current backup
		if len(oldFiles) == len(filenames) {
			if err := os.Rename(oldFiles[len(oldFiles)-1], Filename()); err != nil {
				return err
			}
			oldFiles = oldFiles[:len(oldFiles)-1]
		}

		nod.TotalInt(tpw, len(oldFiles))

		for _, fn := range oldFiles {
			filename := filepath.Join(dir, fn)
			if err := os.Remove(filename); err != nil {
				return err
			}
			nod.Increment(tpw)
		}
	}

	return nil
}
