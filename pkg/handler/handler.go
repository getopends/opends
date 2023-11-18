package handler

import (
	"encoding/json"
	"net/http"
)

type Handler struct{}

func (h Handler) SendJSON(rw http.ResponseWriter, a any) error {
	return json.NewEncoder(rw).Encode(a)
}
