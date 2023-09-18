package server

import (
	"asciify-web/pkg/render"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

// Credit : https://github.com/zenazn/goji/blob/master/web/middleware/nocache.go

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func noCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

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
