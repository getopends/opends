package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/getopends/opends/internal"
	"github.com/gorilla/mux"
)

func main() {
	var (
		fHost    = flag.String("host", "", "Host")
		port     = flag.Int("port", 12345, "Port")
		tls      = flag.Bool("tls", false, "Enable TLS")
		certFile = flag.String("cert-file", "", "Cert file")
		keyfile  = flag.String("key-file", "", "Key file")
	)

	flag.Parse()

	h := &internal.Handler{
		Service:      &internal.Service{},
		PublicRouter: mux.NewRouter(),
		Config: &internal.Config{
			Host: *fHost,
			Port: int16(*port),
		},
	}

	h.SetRoutes()

	host := h.Config.Host

	if host == "" {
		host = "0.0.0.0"
	}

	addr := fmt.Sprintf("%v:%v", host, h.Config.Port)

	srv := http.Server{
		Addr:    addr,
		Handler: h.PublicRouter,
	}

	if *tls {
		addr = fmt.Sprintf("https://%v", addr)
	} else {
		addr = fmt.Sprintf("http://%v", addr)
	}

	log.Printf("Starting server at %v", addr)

	if *tls {
		if err := srv.ListenAndServeTLS(*certFile, *keyfile); err != nil {
			panic(err)
		}
	} else {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}
}

type Response struct {
	Header http.Header
	Body   io.Reader
}
