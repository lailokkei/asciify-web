package server

import (
	"asciify-web/pkg/render"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const staticDir = "static"
const pagesDir = "pages"
const port = "8080"

func pageRerender(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.RenderTemplates()
		h.ServeHTTP(w, r)
	})
}

func catch404(h http.Handler, root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := os.Stat(root + r.URL.Path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				f, err := os.ReadFile(pagesDir + "/404.html")
				if err != nil {
					log.Println(err)
				}
				w.Write([]byte(f))
				return
			}
			log.Println(err)
		}

		h.ServeHTTP(w, r)
	})
}

func Start() {
	var devFlag bool
	flag.BoolVar(&devFlag, "dev", false, "Enables page hot reloading")
	flag.Parse()

	render.RenderTemplates()

	fileServer := catch404(noCache(http.FileServer(http.Dir(staticDir))), staticDir)
	pageServer := noCache(http.FileServer(http.Dir(pagesDir)))

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
