package utils

import (
	"github.com/gorilla/mux"
	"net/http"
)

func QueryString(r *http.Request, s string) string {
	return r.URL.Query().Get(s)
}

func QueryParam(r *http.Request, s string) string {
	params := mux.Vars(r)
	return params[s]
}
