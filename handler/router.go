// handler/router.go
package main

import (
	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/submit", submitHandler).Methods("POST")
	r.HandleFunc("/invoke", invokeHandler).Methods("POST")
	return r
}
