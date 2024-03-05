package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wklken/flow/plugin"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	p := plugin.New("request_id")
	p.Init(`{"header_name": "X-Request-ID", "set_in_response": true}`)

	p1 := plugin.New("basic_auth")
	p1.Init(`{"credentials": {"admin": "admin"}, "realm": "Restricted"}`)

	p2 := plugin.New("file_logger")
	p2.Init(`{"level": "info", "filename": "test.log"}`)

	chain := plugin.BuildPluginChain(p, p1, p2)
	myHandler := http.HandlerFunc(welcomeHandler)

	// ms := make([]func(http.Handler) http.Handler, 2)
	// for _, p := range pChain {
	// 	ms = append(ms, p.Handler)
	// }
	// chi.Chain(...ms).

	r.Handle("/", chain.Then(myHandler))

	// r.Use()
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })

	http.ListenAndServe(":3000", r)
}
