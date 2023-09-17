package server

import (
	"asciify-web/pkg/render"
	"flag"
	"fmt"
	"net/http"
)

const staticDir = "static"
const pagesDir = "pages"
const port = "8080"

func PageRerender(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.RenderTemplates()
		h.ServeHTTP(w, r)
	})
}

func Start() {
	var devFlag bool
	flag.BoolVar(&devFlag, "dev", false, "Enables page hot reloading")
	flag.Parse()

	fs := http.FileServer(http.Dir(staticDir))
	pageServer := NoCache(http.FileServer(http.Dir(pagesDir)))

	if devFlag {
		fmt.Println("--dev")
		pageServer = PageRerender(pageServer)
	}

	http.Handle("/", pageServer)
	http.Handle("/assets/", http.StripPrefix("/assets", NoCache(fs)))
	http.HandleFunc("/connect", connect)

	fmt.Println("serving on : http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
