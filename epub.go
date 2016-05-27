/*
Package epub generates valid EPUB 3.0 files with additional EPUB 2.0 table of
contents (as seen here: https://github.com/bmaupin/epub-samples) for maximum
compatibility.

Basic usage:

	// Create a new EPUB
	e := epub.NewEpub("My title")

	// Set the author
	e.SetAuthor("Hingle McCringleberry")

	// Add a section
	section1Content := `    <h1>Section 1</h1>
	<p>This is a paragraph.</p>`
	e.AddSection("Section 1", section1Content, "section0001.xhtml")

	// The section file name is optional
	section2Content := `    <h1>Section 2</h1>
	<p>This is a paragraph.</p>`
	e.AddSection("Section 2", section2Content, "")

	// Add an image from a local file
	e.AddImage("testdata/gophercolor16x16.png", "go-gopher.png")

	// Add an image from a URL. The image filename is also optional
	e.AddImage("https://golang.org/doc/gopher/gophercolor16x16.png", "")

	// Write the EPUB
	err = e.Write("My EPUB.epub")
	if err != nil {
		// handle error
	}

*/
package epub

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/satori/go.uuid"
)

// ErrFilenameAlreadyUsed is thrown by AddImage and AddSection if the same
// filename is used more than once
var ErrFilenameAlreadyUsed = errors.New("Filename already used")

const (
	defaultEpubLang   = "en"
	imageFileFormat   = "image%04d%s"
	sectionFileFormat = "section%04d.xhtml"
	urnUUIDPrefix     = "urn:uuid:"
)

// Epub implements an EPUB file.
type Epub struct {
	author   string
	images   map[string]string // Images added to the EPUB
	lang     string            // Language
	pkg      *pkg              // The package file (package.opf)
	sections map[string]xhtml  // Sections (chapters)
	title    string
	toc      *toc // Table of contents
	uuid     string
}

// NewEpub returns a new Epub.
func NewEpub(title string) *Epub {
	e := &Epub{}
	e.images = make(map[string]string)
	e.sections = make(map[string]xhtml)
	e.pkg = newPackage()
	e.toc = newToc()
	// Set minimal required attributes
	e.SetLang(defaultEpubLang)
	e.SetTitle(title)
	e.SetUUID(uuid.NewV4().String())

	return e
}

// AddImage adds an image to the EPUB and returns a relative path that can be
// used in the content of a section.
//
// The image source should either be a URL or a path to a local file; in either
// case, the image will be retrieved and stored in the EPUB.
//
// The image filename will be used when storing the image in the EPUB and must
// be unique. If the same filename is used more than once,
// ErrFilenameAlreadyUsed will be returned. The image filename is optional; if
// no filename is provided, one will be generated.
func (e *Epub) AddImage(imageSource string, imageFilename string) (string, error) {
	// Generate an image filename if one isn't provided
	if imageFilename == "" {
		imageFilename = fmt.Sprintf(imageFileFormat, len(e.images)+1, filepath.Ext(imageSource))
	}

	if _, ok := e.images[imageFilename]; ok {
		return "", ErrFilenameAlreadyUsed
	}

	e.images[imageFilename] = imageSource

	return filepath.Join(
		"..",
		imageFolderName,
		imageFilename,
	), nil
}

// AddSection adds a new section (chapter, etc) to the EPUB and returns a
// relative path to the section that can be used from another section (for
// links).
//
// The title will be used for the table of contents.
//
// The content must be valid XHTML that will go between the <body> tags of the
// section XHTML file. The content will not be validated.
//
// The section filename will be used when storing the image in the EPUB and must
// be unique. If the same filename is used more than once,
// ErrFilenameAlreadyUsed will be returned. The section filename is optional;
// if no filename is provided, one will be generated.
//
// The section will be shown in the table of contents in the same order it was
// added to the EPUB.
func (e *Epub) AddSection(sectionTitle string, sectionContent string, sectionFilename string) (string, error) {
	// Generate a section filename if one isn't provided
	if sectionFilename == "" {
		sectionFilename = fmt.Sprintf(sectionFileFormat, len(e.sections)+1)
	}

	if _, ok := e.sections[sectionFilename]; ok {
		return "", ErrFilenameAlreadyUsed
	}

	x := newXhtml(sectionContent)
	x.setTitle(sectionTitle)

	e.sections[sectionFilename] = *x

	return sectionFilename, nil
}

// Author returns the author of the EPUB.
func (e *Epub) Author() string {
	return e.author
}

// Lang returns the language of the EPUB.
func (e *Epub) Lang() string {
	return e.lang
}

// SetAuthor sets the author of the EPUB.
func (e *Epub) SetAuthor(author string) {
	e.author = author
	e.pkg.setAuthor(author)
}

// SetLang sets the language of the EPUB.
func (e *Epub) SetLang(lang string) {
	e.lang = lang
	e.pkg.setLang(lang)
}

// SetTitle sets the title of the EPUB.
func (e *Epub) SetTitle(title string) {
	e.title = title
	e.pkg.setTitle(title)
	e.toc.setTitle(title)
}

// SetUUID sets the UUID of the EPUB. A UUID will be automatically be generated
// for you when the NewEpub method is run.
func (e *Epub) SetUUID(uuid string) {
	e.uuid = uuid
	e.pkg.setUUID(urnUUIDPrefix + uuid)
	e.toc.setUUID(urnUUIDPrefix + uuid)
}

// Title returns the title of the EPUB.
func (e *Epub) Title() string {
	return e.title
}

// UUID returns the UUID of the EPUB.
func (e *Epub) UUID() string {
	return e.uuid
}
