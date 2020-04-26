package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("wellcome"))
	})
	http.ListenAndServe(":3000", r)
}
