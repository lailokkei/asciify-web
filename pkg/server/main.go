package server

import (
	"asciify-web/pkg/render"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

const staticDir = "static"
const pagesDir = "pages"
const port = "8080"

func Start() {
	log.SetFlags(0)
	log.SetOutput(io.Discard)
	var devFlag bool
	flag.BoolVar(&devFlag, "dev", false, "Enables page hot reloading")
	flag.Parse()

	render.RenderTemplates()

	fileServer := noCache(http.FileServer(http.Dir(staticDir)))
	pageServer := catch404(noCache(http.FileServer(http.Dir(pagesDir))), pagesDir)

	if devFlag {
		fmt.Println("--dev")
		pageServer = pageRerender(pageServer)
		fileServer = pageRerender(fileServer)
	}

	http.Handle("/", pageServer)
	http.Handle("/assets/", http.StripPrefix("/assets", fileServer))
	http.HandleFunc("/connect", connect)

	fmt.Println("serving on : http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
