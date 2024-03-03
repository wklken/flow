package main

import (
	"net/http"
	"sort"

	"github.com/justinas/alice"

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
	p := plugin.Plugins["request_id"]
	p.Init(`{"header_name": "X-Request-ID", "set_in_response": true}`)
	// p.Init(`{}`)

	p1 := plugin.Plugins["basic_auth"]
	p1.Init(`{"credentials": {"admin": "admin"}, "realm": "Restricted"}`)

	pChain := []plugin.Plugin{p, p1}
	sort.Slice(pChain, func(i, j int) bool {
		return pChain[i].Priority() < pChain[j].Priority()
	})
	chain := alice.New()
	for _, p := range pChain {
		chain = chain.Append(p.Handler)
	}
	myHandler := http.HandlerFunc(welcomeHandler)

	r.Handle("/", chain.Then(myHandler))

	// r.Use()
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })

	http.ListenAndServe(":3000", r)
}
