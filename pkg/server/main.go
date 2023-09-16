package server

import (
	"asciify-web/pkg/render"
	"fmt"
	"net/http"
)

// func catch404(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
//     fn := http
//     return
// }

const staticDir = "static"
const pagesDir = "pages"
const port = "8080"

func Start() {
	render.RenderTemplates()

	fs := http.FileServer(http.Dir(staticDir))
	ps := http.FileServer(http.Dir(pagesDir))

	http.Handle("/", NoCache(ps))
	http.Handle("/assets/", http.StripPrefix("/assets", NoCache(fs)))
	http.HandleFunc("/connect", connect)

	fmt.Println("serving on : http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
