package routes

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n", r.Header.Get("Authorization"))
	w.WriteHeader(204)
}

