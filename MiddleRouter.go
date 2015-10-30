package main

import (
    "net/http"

    "github.com/gorilla/mux"
)

type MiddleRouter struct {
	mux *mux.Router
}

func (s *MiddleRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if origin := r.Header.Get("Origin"); origin != "" {
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    }
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }

    w.Header().Set("Master", "JRod")

    s.mux.ServeHTTP(w, r) // Lets Gorilla work
}
