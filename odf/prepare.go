package odf

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
)

var errNotImplemented = errors.New("odf: not impemented")

// OpenDocument represents an open document reader and writter.
type OpenDocument struct {
	zipFile *os.File
	cache   map[string]*bytes.Buffer
}

// Open loads the file contents in-memory
func Open(filename string) (*OpenDocument, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	odf := &OpenDocument{}
	odf.initCache()

	for _, f := range r.File {
		logger.Printf("archive=%v; file=%v; dir=%v", filename, f.Name, f.FileInfo().IsDir())
		if f.FileInfo().IsDir() {
			continue
		}
		buff := odf.fileBuffer(f.Name)
		zipfd, err := f.Open()
		if err != nil {
			logger.Printf("Error opening file %v from zip: %v", f, err)
			continue
		}
		w, err := io.Copy(buff, zipfd)
		if err != nil {
			logger.Printf("Error reading file %v from zip: %v", f, err)
			continue
		}
		if uint64(w) != f.UncompressedSize64 {
			logger.Printf("Error reading file: not all bytes read, expected %d, got %d", f.UncompressedSize64, w)
		}
		logger.Printf("Read %d bytes from %v", w, f.Name)
	}
	return odf, nil
}

// Prepare parses annotations in the Open Document and interpolates them
// on the content.xml and styles.xml files, allowing users to easily "prepare"
// reports for future rendering.
func (o *OpenDocument) Prepare() error {
	return errNotImplemented
}

// Save package and write the file contents into the destination.
func (o *OpenDocument) Save(filename string) error {
	return errNotImplemented
}

// Extract decodes the file and extracts it's content into the specified folder
func (o *OpenDocument) Extract(folder string) error {
	return errNotImplemented
}

func (o *OpenDocument) initCache() {
	o.cache = make(map[string]*bytes.Buffer)
}

func (o *OpenDocument) listFiles() (files []string) {
	for k := range o.cache {
		files = append(files, k)
	}
	return
}

func (o *OpenDocument) fileBuffer(filename string) *bytes.Buffer {
	if b, ok := o.cache[filename]; ok {
		logger.Printf("Buffer for %v found: %p", filename, b)
		return b
	}
	b := &bytes.Buffer{}
	logger.Printf("New buffer for %v created: %p", filename, b)
	o.cache[filename] = b
	return b
}
