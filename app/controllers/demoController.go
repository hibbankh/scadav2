package controllers

import (
	"net/http"

	"app/network"

	"github.com/go-chi/chi"
)

type (
	DemoResponse struct {
		Id   uint64
		Data string
	}
)

func DemoRoute() *chi.Mux {
	router := chi.NewMux()
	router.Get("/", index)
	router.Post("/", store)
	router.Get("/{id}", show)
	router.Put("/{id}", update)
	router.Patch("/{id}", update)
	router.Delete("/{id}", delete)

	return router
}

func index(w http.ResponseWriter, r *http.Request) {
	var response DemoResponse
	network.ResponseJSON(
		w, false, http.StatusOK, response,
	)
}

func store(w http.ResponseWriter, r *http.Request) {}

func show(w http.ResponseWriter, r *http.Request) {}

func update(w http.ResponseWriter, r *http.Request) {}

func delete(w http.ResponseWriter, r *http.Request) {}
