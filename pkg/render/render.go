package render

import (
	"os"
	"text/template"
)

type page struct {
	Header  string
	Content string
}

const templatesDir = "templates/"
const pagesDir = templatesDir + "pages/"
const outputDir = "pages/"

const layoutFile = "layout.tmpl"
const headerFile = templatesDir + "header.html"

func RenderTemplates() {
	tmpl, err := template.New(layoutFile).ParseFiles(templatesDir + layoutFile)
	if err != nil {
		panic(err)
	}

	header, err := os.ReadFile(headerFile)
	if err != nil {
		panic(err)
	}

	pages, err := os.ReadDir(pagesDir)
	if err != nil {
		panic(err)
	}

	for _, p := range pages {
		file, err := os.Create(outputDir + p.Name())
		if err != nil {
			panic(err)
		}

		content, err := os.ReadFile(pagesDir + p.Name())
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(file, page{string(header), string(content)})
		if err != nil {
			panic(err)
		}
	}
}
