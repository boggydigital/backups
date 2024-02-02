package hogo

import (
	"github.com/boggydigital/nod"
	"time"
)

const (
	tarGzExt = ".tar.gz"
)

func Filename() string {
	return time.Now().Format(nod.TimeFormat) + tarGzExt
}
