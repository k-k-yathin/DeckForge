package service

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"strings"
)

// readDocxText extracts visible text from a .docx file without external deps.
// DOCX files are ZIP archives containing word/document.xml.
func readDocxText(path string) (string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		return parseWordXML(rc)
	}
	return "", nil
}

// parseWordXML walks WordprocessingML and collects text from <w:t> nodes.
func parseWordXML(r io.Reader) (string, error) {
	decoder := xml.NewDecoder(r)
	var parts []string
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			if el.Name.Local == "t" {
				var text string
				if err := decoder.DecodeElement(&text, &el); err == nil && text != "" {
					parts = append(parts, text)
				}
			}
		}
	}
	return strings.Join(parts, " "), nil
}
