package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	// Packages
	pg "github.com/mutablelogic/go-pg"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func RegisterHandlers(router *http.ServeMux, conn pg.Conn) http.Handler {
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Get the id value
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Call handler
		switch r.Method {
		case http.MethodGet:
			GetHandler(w, r, conn, id)
		case http.MethodPatch:
			PatchHandler(w, r, conn, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		// Call handler
		switch r.Method {
		//case http.MethodPost:
		//	PostHandler(w, r, conn)
		case http.MethodGet:
			ListHandler(w, r, conn)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return router
}

// GET /
func ListHandler(w http.ResponseWriter, r *http.Request, conn pg.Conn) {
	var response NameList
	if err := conn.List(r.Context(), &response, NameList{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(response)
}

// GET /{id}
func GetHandler(w http.ResponseWriter, r *http.Request, conn pg.Conn, id uint64) {
	var response Name
	if err := conn.Get(r.Context(), &response, Name{Id: id}); errors.Is(err, pg.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// PATCH /{id}
func PatchHandler(w http.ResponseWriter, r *http.Request, conn pg.Conn, id uint64) {
	var request Name
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var response Name
	if err := conn.Update(r.Context(), &response, request, request); errors.Is(err, pg.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
