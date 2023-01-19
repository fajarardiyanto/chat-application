package utils

import "net/http"

func QueryString(r *http.Request, s string) string {
	return r.URL.Query().Get(s)
}
