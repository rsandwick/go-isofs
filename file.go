package isofs

import (
	"io"
	"os"

	"github.com/kdomanski/iso9660"
)

type file struct {
	fi     *iso9660.File
	r      io.Reader
	closed bool
}

func (f *file) Close() error {
	f.closed = true
	return nil
}

func (f *file) Read(buf []byte) (int, error) {
	if f == nil || f.closed {
		return 0, os.ErrInvalid
	}
	return f.r.Read(buf)
}

func (f *file) Stat() (os.FileInfo, error) {
	if f == nil || f.closed {
		return nil, os.ErrInvalid
	}
	return f.fi, nil
}
