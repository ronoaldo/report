package odf

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/beevik/etree"
)

var errNotImplemented = errors.New("odf: not impemented")

// OpenDocument represents an open document reader and writter.
type OpenDocument struct {
	zipFile *os.File
	cache   map[string]*bytes.Buffer
}

// Open opens, reads and unpacks the document located at filename in memory.
func Open(filename string) (*OpenDocument, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return OpenReader(r)
}

// OpenReader reads and unpacks the provided document from the provided zip reader in memory.
func OpenReader(r *zip.ReadCloser) (*OpenDocument, error) {
	odf := &OpenDocument{}
	odf.initCache()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		buff := odf.fileBuffer(f.Name)
		zipfd, err := f.Open()
		if err != nil {
			printf("Error opening file %v from zip: %v", f, err)
			continue
		}
		w, err := io.Copy(buff, zipfd)
		if err != nil {
			printf("Error reading file %v from zip: %v", f, err)
			continue
		}
		if uint64(w) != f.UncompressedSize64 {
			printf("Error reading file: not all bytes read, expected %d, got %d", f.UncompressedSize64, w)
		}
		printf("Read %d bytes from %v", w, f.Name)
	}
	return odf, nil
}

// Execute merges the provided data into all loaded and parsed template files.
func (o *OpenDocument) Execute(values interface{}) error {
	filesToMerge := []string{"content.xml", "styles.xml"}
	for _, fname := range filesToMerge {
		f := o.fileBuffer(fname)
		xml := f.String()
		if fname == "content.xml" {
			var err error
			xml, err = prepareXMLForTemplate(xml)
			if err != nil {
				return err
			}
		}
		xml = formatXML(xml)
		printf("PARSING: %s -> %v", fname, xml)
		tpl, err := template.New(fname).Parse(xml)
		if err != nil {
			return fmt.Errorf("odf: unable to parse template: %v [%v]", err, xml)
		}
		tpl = tpl.Option("missingkey=zero")
		buff := &bytes.Buffer{}
		if err = tpl.Execute(buff, values); err != nil {
			return err
		}
		f.Reset()
		b, err := io.Copy(f, buff)
		if err != nil {
			return err
		}
		printf("Merged %s (%d bytes)", fname, b)
	}
	return nil
}

// WriteFile package and write the file contents into the destination.
func (o *OpenDocument) WriteFile(filename string) error {
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	w := zip.NewWriter(fd)
	defer fd.Close()

	// We need to overwrite each zip file
	for k, buff := range o.cache {
		zipfd, err := w.Create(k)
		if err != nil {
			return err
		}
		b, err := io.Copy(zipfd, buff)
		if err != nil {
			return err
		}
		printf("Written %d bytes for %s into zip.", b, k)
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return w.Close()
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
		return b
	}
	b := &bytes.Buffer{}
	o.cache[filename] = b
	return b
}

func prepareXMLForTemplate(rawXML string) (string, error) {
	// Parse XML file using DOM
	doc := etree.NewDocument()
	err := doc.ReadFromString(rawXML)
	if err != nil {
		return rawXML, err
	}

	// Search for paragraphs and fix up invalid elements
	// TODO: optimize the processing here to allow more robust
	// templating.
	toClean := []string{"//text:p", "//annotation"}
	for _, path := range toClean {
		for _, p := range doc.FindElements(path) {
			var prev *etree.Element
			spans := p.ChildElements()
			for i := range spans {
				s := spans[i]
				if s.Tag != "span" {
					prev = nil
					continue
				}
				if s.Tag == "s" {
					printf("Removed <text:s/>")
					p.RemoveChild(s)
					prev = nil
					continue
				}
				// Let's clean up some dup styles that libreoffice goes on crazy building.
				if prev == nil {
					prev = s
					continue
				}

				if s.SelectAttr("style-name").Value == prev.SelectAttr("style-name").Value {
					printf("Found matching style from previous tag. merge in the contents")
					prev.SetText(prev.Text() + s.Text())
					p.RemoveChild(s)
				}
			}
		}
	}

	// Remove annotations, and move the contents to up/bottom of table row
	for _, a := range doc.FindElements("//annotation") {
		content := ""
		for _, e := range a.ChildElements() {
			if e.Tag == "p" {
				span := e.SelectElement("span")
				if span == nil {
					continue
				}
				content = span.Text()
			}
		}
		if strings.Contains(content, "range") || strings.Contains(content, "end") {
			tr := a.Parent()
			for tr.Tag != "body" && tr.Tag != "table-row" {
				tr = tr.Parent()
			}
			if tr.Tag == "body" {
				return rawXML, fmt.Errorf("odf: found annotation with %s but no parent table-row found", content)
			}
			table := tr.Parent()
			textNode := etree.NewCharData(content)

			if strings.Contains(content, "range") {
				table.InsertChild(tr, textNode)
			} else {
				next := nextSimbling(tr)
				if next == tr {
					table.AddChild(textNode)
				} else {
					table.InsertChild(nextSimbling(tr), textNode)
				}
			}
		}
		a.Parent().RemoveChild(a)
	}

	// Render new XML
	cleanXML, err := doc.WriteToString()
	if err != nil {
		return rawXML, err
	}
	return cleanXML, err
}

func formatXML(src string) string {
	// Parse XML file using DOM
	doc := etree.NewDocument()
	err := doc.ReadFromString(src)
	if err != nil {
		return src
	}
	doc.Indent(2)
	xml, err := doc.WriteToString()
	if err != nil {
		return src
	}
	// Quick and dirty fix to allow us to use '""' inside template directives
	return strings.Replace(xml, "&quot;", "\"", -1)
}

func nextSimbling(e *etree.Element) *etree.Element {
	if e == nil {
		return e
	}
	printf("> Lookup for next simbling for <%s:%s>", e.Space, e.Tag)
	parent := e.Parent()
	if parent == nil {
		// No parent, next simbling is itself
		return e
	}
	printf("> Found parent <%s:%s>", parent.Space, parent.Tag)
	children := parent.ChildElements()
	if children == nil {
		return e
	}
	pos := -1
	for i := range children {
		c := children[i]
		if c == e {
			pos = i
			break
		}
	}
	if pos == -1 {
		panic(fmt.Sprintf("odf: cannot find next simbling of %v", e))
	}
	if pos >= len(children)-1 {
		// last element
		return e
	}
	return children[pos+1]
}
