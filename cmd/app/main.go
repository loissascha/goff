package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/loissascha/goff/pkg/framework"
)

func main() {
	addr := env("ADDR", ":8118")

	renderer, err := framework.NewRenderer(framework.RendererConfig{
		LayoutGlob:   "internal/pages/layouts/*.html",
		TemplateGlob: "internal/pages/templates/*.html",
	})
	if err != nil {
		log.Fatalf("load templates: %v", err)
	}

	router := framework.NewRouter(renderer)
	router.Page(http.MethodGet, "/", homePage())

	mux := http.NewServeMux()
	mux.Handle("/", router)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func env(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func homePage() framework.Handler {
	return func(ctx *framework.Context) (framework.Page, error) {
		return framework.Page{
			Title:       "Test Page",
			Description: "Test Page Description",
			Template:    "home.html",
		}, nil
	}
}
