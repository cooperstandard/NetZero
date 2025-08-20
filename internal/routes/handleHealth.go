package routes

import "net/http"

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
