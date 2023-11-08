package main

import (
	"context"
	"fmt"
	"github.com/ossrs/go-oryx-lib/logger"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	ctx := logger.WithContext(context.Background())

	setDefaultEnv := func(k, v string) {
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
	setDefaultEnv("LISTEN", ":2025")

	staticHandler := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".js") {
			// All js files are in js folder.
			r.URL.Path = path.Join("js", path.Base(r.URL.Path))

			staticHandler.ServeHTTP(w, r)
			return
		}

		if r.URL.RawQuery == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "no-cache")

		logger.Tf(ctx, "Handle %v", r.URL.RequestURI())

		tmpl := template.New("index.tmpl.html")
		tmpl = template.Must(tmpl.ParseFiles("index.tmpl.html"))
		tmpl.Execute(w, &struct {
			Target string
		}{
			Target: r.URL.Path,
		})
	})

	addr := os.Getenv("LISTEN")
	if !strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf(":%v", os.Getenv("LISTEN"))
	}
	logger.Tf(ctx, "Listen on %v", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Tf(ctx, "%v", err)
	}
}
