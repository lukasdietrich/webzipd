package render

import (
	"bytes"
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
		contentReader := writeContentType(hw, reader)

		io.Copy(hw, contentReader)
		reader.Close()
	}
}

func writeContentType(hw http.ResponseWriter, reader *namespace.Reader) io.Reader {
	contentReader := io.Reader(reader)
	contentType := mime.TypeByExtension(filepath.Ext(reader.Filename))

	if contentType == "" {
		guessBuffer := make([]byte, 512)
		n, _ := io.ReadFull(reader, guessBuffer)
		contentType = http.DetectContentType(guessBuffer[:n])
		contentReader = io.MultiReader(bytes.NewReader(guessBuffer[:n]), contentReader)
	}

	hw.Header().Set("Content-Type", contentType)

	return contentReader
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
