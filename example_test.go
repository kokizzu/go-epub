package epub_test

import (
	"fmt"
	"log"

	"github.com/bmaupin/go-epub"
)

func ExampleEpub_AddCSS() {
	e := epub.NewEpub("My title")

	cssContent := `h1 {
  text-align: center;
}`

	// Add CSS
	css1Path, err := e.AddCSS(cssContent, "epub.css")
	if err != nil {
		log.Fatal(err)
	}

	// The filename is optional
	css2Path, err := e.AddCSS(cssContent, "")
	if err != nil {
		log.Fatal(err)
	}

	// Use the CSS in a section
	sectionContent := `    <h1>Section 1</h1>
	<p>This is a paragraph.</p>`
	e.AddSection("Section 1", sectionContent, "", css1Path)

	fmt.Println(css1Path)
	fmt.Println(css2Path)

	// Output:
	// ../css/epub.css
	// ../css/css0002.css
}

func ExampleEpub_AddImage() {
	e := epub.NewEpub("My title")

	// Add an image from a local file
	img1Path, err := e.AddImage("testdata/gophercolor16x16.png", "go-gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	// Add an image from a URL. The filename is optional
	img2Path, err := e.AddImage("https://golang.org/doc/gopher/gophercolor16x16.png", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(img1Path)
	fmt.Println(img2Path)

	// Output:
	// ../img/go-gopher.png
	// ../img/image0002.png
}

func ExampleEpub_AddSection() {
	e := epub.NewEpub("My title")

	// Add a section
	section1Content := `    <h1>Section 1</h1>
	<p>This is a paragraph.</p>`
	section1Path, err := e.AddSection("Section 1", section1Content, "firstsection.xhtml", "")
	if err != nil {
		log.Fatal(err)
	}

	// Link to the first section
	section2Content := fmt.Sprintf(`    <h1>Section 2</h1>
	<a href="%s">Link to section 1</a>`,
		section1Path)
	// The filename is optional
	section2Path, err := e.AddSection("Section 2", section2Content, "", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(section1Path)
	fmt.Println(section2Path)

	// Output:
	// firstsection.xhtml
	// section0002.xhtml
}