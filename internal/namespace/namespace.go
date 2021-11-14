package namespace

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Reader struct {
	io.ReadCloser
	Filename string
	Modtime  time.Time
}

type Namespace struct {
	mu       sync.Mutex
	modtime  time.Time
	filename string

	zip        *zip.ReadCloser
	contentMap map[string]*zip.File
}

func openNamespace(filename string) (*Namespace, error) {
	n := &Namespace{
		filename: filename,
	}

	return n, n.load()
}

func (n *Namespace) load() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	info, err := os.Stat(n.filename)
	if err != nil {
		return err
	}

	if !info.ModTime().After(n.modtime) {
		return nil
	}

	if n.zip != nil {
		if err := n.zip.Close(); err != nil {
			return err
		}
	}

	n.zip, err = zip.OpenReader(n.filename)
	if err != nil {
		return err
	}

	n.contentMap = make(map[string]*zip.File, len(n.zip.File))
	for _, file := range n.zip.File {
		n.contentMap[file.Name] = file
	}

	return nil
}

func (n *Namespace) Open(filename string) (*Reader, error) {
	if err := n.load(); err != nil {
		return nil, err
	}

	filename = strings.TrimLeft(filename, "/")

	file, ok := n.contentMap[filename]
	if !ok || file.Mode().IsDir() {
		if !strings.HasSuffix(filename, "index.html") {
			return n.Open(path.Join(filename, "index.html"))
		}

		return nil, os.ErrNotExist
	}

	r, err := file.Open()
	return &Reader{ReadCloser: r, Filename: filename, Modtime: file.Modified}, err
}
