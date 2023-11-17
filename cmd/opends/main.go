package main

import (
	"io"
	"net/http"

	"github.com/getopends/opends/internal"
	"github.com/gorilla/mux"
)

func main() {
	h := &internal.Handler{
		Service:      &internal.Service{},
		PublicRouter: mux.NewRouter(),
	}

	h.SetPublicRoutes()

	srv := http.Server{
		Addr:    ":12345",
		Handler: h.PublicRouter,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

type Response struct {
	Header http.Header
	Body   io.Reader
}
