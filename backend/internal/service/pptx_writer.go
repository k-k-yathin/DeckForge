package service

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/deckforge/backend/internal/models"
)

// writePPTX builds a valid .pptx file (Office Open XML ZIP package).
// This avoids heavy dependencies and works in Docker without CGO.
func writePPTX(p *models.Presentation, outPath string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()
	return writePPTXZip(zw, p)
}

func writePPTXZip(z *zip.Writer, p *models.Presentation) error {
	slideCount := len(p.Slides)
	if slideCount == 0 {
		slideCount = 1
	}

	// Package relationships
	if err := writeZipFile(z, "_rels/.rels", rootRels()); err != nil {
		return err
	}
	if err := writeZipFile(z, "[Content_Types].xml", contentTypes(slideCount)); err != nil {
		return err
	}
	if err := writeZipFile(z, "ppt/presentation.xml", presentationXML(slideCount)); err != nil {
		return err
	}
	if err := writeZipFile(z, "ppt/_rels/presentation.xml.rels", presentationRels(slideCount)); err != nil {
		return err
	}

	for i, slide := range p.Slides {
		n := i + 1
		var bullets []string
		_ = json.Unmarshal(slide.Content, &bullets)
		body := buildSlideBody(slide.Title, slide.Subtitle, bullets)
		if err := writeZipFile(z, fmt.Sprintf("ppt/slides/slide%d.xml", n), slideXML(body)); err != nil {
			return err
		}
		if err := writeZipFile(z, fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", n), slideRels()); err != nil {
			return err
		}
	}

	return nil
}

func writeZipFile(z *zip.Writer, name, content string) error {
	w, err := z.Create(name)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, content)
	return err
}

func buildSlideBody(title, subtitle string, bullets []string) string {
	var b strings.Builder
	b.WriteString(escapeXML(title))
	if subtitle != "" {
		b.WriteString("\n")
		b.WriteString(escapeXML(subtitle))
	}
	for _, bullet := range bullets {
		b.WriteString("\n• ")
		b.WriteString(escapeXML(bullet))
	}
	return b.String()
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

func contentTypes(slideCount int) string {
	var overrides strings.Builder
	overrides.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
`)
	for i := 1; i <= slideCount; i++ {
		overrides.WriteString(fmt.Sprintf(`  <Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`+"\n", i))
	}
	overrides.WriteString(`</Types>`)
	return overrides.String()
}

func rootRels() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
</Relationships>`
}

func presentationRels(slideCount int) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for i := 1; i <= slideCount; i++ {
		b.WriteString(fmt.Sprintf(`
  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide%d.xml"/>`, i, i))
	}
	b.WriteString("\n</Relationships>")
	return b.String()
}

func presentationXML(slideCount int) string {
	var ids strings.Builder
	for i := 1; i <= slideCount; i++ {
		ids.WriteString(fmt.Sprintf(`<p:sldId id="%d" r:id="rId%d"/>`, 256+i, i))
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:sldIdLst>%s</p:sldIdLst>
  <p:sldSz cx="9144000" cy="6858000"/>
</p:presentation>`, ids.String())
}

func slideRels() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>
`
}

func slideXML(text string) string {
	// Split text into lines for separate paragraphs
	lines := strings.Split(text, "\n")
	var paras strings.Builder
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		paras.WriteString(fmt.Sprintf(`<a:p><a:r><a:t>%s</a:t></a:r></a:p>`, line))
	}
	if paras.Len() == 0 {
		paras.WriteString(`<a:p><a:r><a:t> </a:t></a:r></a:p>`)
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <p:cSld>
    <p:spTree>
      <p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>
      <p:grpSpPr/>
      <p:sp>
        <p:nvSpPr><p:cNvPr id="2" name="Content"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr/>
        <p:txBody>
          <a:bodyPr/>
          <a:lstStyle/>
          %s
        </p:txBody>
      </p:sp>
    </p:spTree>
  </p:cSld>
</p:sld>`, paras.String())
}
