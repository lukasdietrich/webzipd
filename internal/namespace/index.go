package namespace

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Index struct {
	mu         sync.Mutex
	foldername string

	namespaceMap map[string]*Namespace
}

func OpenIndex(foldername string) (*Index, error) {
	i := &Index{
		foldername: foldername,
	}

	return i, i.load()
}

func (i *Index) load() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	folder, err := os.Open(i.foldername)
	if err != nil {
		return fmt.Errorf("could not open folder %q: %w", i.foldername, err)
	}

	defer folder.Close()

	infoSlice, err := folder.Readdir(-1)
	if err != nil {
		return fmt.Errorf("could not list files in folder %q: %w", i.foldername, err)
	}

	i.namespaceMap = make(map[string]*Namespace, len(infoSlice))

	for _, info := range infoSlice {
		if name := info.Name(); !info.IsDir() && strings.HasSuffix(name, ".zip") {
			log.Printf("load namespace %q.", name)

			filename := filepath.Join(i.foldername, name)
			namespace, err := openNamespace(filename)
			if err != nil {
				return fmt.Errorf("could not read content zip %q: %w", name, err)
			}

			i.namespaceMap[name[:len(name)-4]] = namespace
		}
	}

	return nil
}

func (i *Index) Open(namespace, filename string) (*Reader, error) {
	n, ok := i.namespaceMap[namespace]
	if !ok {
		return nil, os.ErrNotExist
	}

	return n.Open(filename)
}
