package render

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/lukasdietrich/webzipd/internal/namespace"
)

type Renderer struct {
	index *namespace.Index
}

func NewRenderer(index *namespace.Index) *Renderer {
	return &Renderer{index: index}
}

func (r *Renderer) Render(hw http.ResponseWriter, hr *http.Request, namespace string) {
	filename := hr.URL.Path

	reader, err := r.index.Open(namespace, filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeStatus(hw, 404)
		} else {
			writeStatus(hw, 500)
		}
	} else {
		writeModTime(hw, reader.Modtime)
		writeContentType(hw, filename)

		io.Copy(hw, reader)
		reader.Close()
	}
}

func writeContentType(hw http.ResponseWriter, filename string) {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	hw.Header().Set("Content-Type", contentType)
}

func writeModTime(hw http.ResponseWriter, modtime time.Time) {
	// See <https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Last-Modified>
	const headerTimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
	hw.Header().Set("Last-Modified", modtime.UTC().Format(headerTimeFormat))
}

func writeStatus(hw http.ResponseWriter, status int) {
	hw.WriteHeader(status)
	hw.Write([]byte(http.StatusText(status)))
}
