package handler

import (
	"encoding/json"
	"net/http"
)

type Base struct{}

func (b Base) SendJSON(rw http.ResponseWriter, a any) error {
	return json.NewEncoder(rw).Encode(a)
}
