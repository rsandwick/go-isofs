package isofs

import (
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"

	bufra "github.com/avvmoto/buf-readerat"
	"github.com/kdomanski/iso9660"
	httpra "github.com/snabb/httpreaderat"
)

type FS struct {
	io.Closer
	*iso9660.Image
	dirs map[string]*dir
}

var _ fs.FS = &FS{}

func Open(x string) (*FS, error) {
	f, err := open(x)
	if err != nil {
		return nil, err
	}
	i, err := iso9660.OpenImage(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	iso := &FS{
		Closer: f,
		Image:  i,
		dirs:   make(map[string]*dir),
	}
	return iso, nil
}

func (iso *FS) Open(path string) (fs.File, error) {
	return &file{}, nil
}

type readAtCloser interface {
	io.Closer
	io.ReaderAt
}

type nopCloser struct {
	io.ReaderAt
}

func (nopCloser) Close() error { return nil }

func open(rawURL string) (readAtCloser, error) {
	u, _ := url.Parse(rawURL)
	switch u.Scheme {
	case "http", "https":
		return urlopen(u)
	case "":
		return os.Open(u.Path)
	default:
		return nil, os.ErrInvalid
	}
}

func urlopen(u *url.URL) (readAtCloser, error) {
	req, _ := http.NewRequest("GET", u.String(), nil)
	hra, err := httpra.New(nil, req, nil)
	if err != nil {
		return nil, err
	}
	return nopCloser{bufra.NewBufReaderAt(hra, 1<<20)}, nil
}

type dir struct {
	entries  []*file
	transTbl map[string]int
}
