package stm

import (
	"bytes"
	"log"
)

// builderFileError is implementation for the BuilderError interface.
type builderFileError struct {
	error
	full bool
}

// FullError returns true if a sitemap xml had been limit size.
func (e *builderFileError) FullError() bool {
	return e.full
}

// NewBuilderFile returns the created the BuilderFile's pointer
func NewBuilderFile(sm *Sitemap) *BuilderFile {
	loc := sm.opts.Location()
	b := &BuilderFile{loc: loc, sitemap: sm}
	b.clear()
	return b
}

// BuilderFile provides implementation for the Builder interface.
type BuilderFile struct {
	content []byte
	loc     *Location
	linkcnt int
	newscnt int
	sitemap *Sitemap
}

// Add method joins old bytes with creates bytes by it calls from Sitemap.Add method.
func (b *BuilderFile) Add(url interface{}) BuilderError {
	u := MergeMap(url.(URL),
		URL{"host": b.loc.opts.defaultHost},
	)

	smu, err := NewSitemapURL(u)
	if err != nil {
		log.Fatalf("[F] Sitemap: %s", err)
	}

	bytes := smu.XML()

	if !b.isFileCanFit(bytes) {
		return &builderFileError{error: err, full: true}
	}

	b.content = append(b.content, bytes...)
	b.linkcnt++
	return nil
}

// isFileCanFit checks bytes to bigger than consts values.
func (b *BuilderFile) isFileCanFit(bytes []byte) bool {
	r := len(append(b.content, bytes...)) < MaxSitemapFilesize
	r = r && b.linkcnt < MaxSitemapLinks
	return r && b.newscnt < MaxSitemapNews
}

// clear will initialize xml content.
func (b *BuilderFile) clear() {
	b.content = make([]byte, 0, MaxSitemapFilesize)
}

// Content will return pooled bytes on content attribute.
func (b *BuilderFile) Content() []byte {
	return b.content
}

// Write will write pooled bytes with header and footer to
// Location path for output sitemap file.
func (b *BuilderFile) Write() {
	b.loc.ReserveName()
	xmlHeader := XMLHeader

	//百度的比较特殊
	if b.sitemap.opts.SearchEngine == SearcnEngine_Baidu {
		xmlHeader = Baidu_XMLHeader
	}

	c := bytes.Join(bytes.Fields(xmlHeader), []byte(" "))
	c = append(append(c, b.Content()...), XMLFooter...)

	b.loc.Write(c, b.linkcnt)
	b.clear()
}
