package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/getopends/opends/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var defaultPort = 13000

func serveCmd(cfg *internal.Config) *cobra.Command {
	serveCmd := &cobra.Command{
		Use: "serve",
		RunE: func(_ *cobra.Command, _ []string) error {
			return doServe(cfg)
		},
	}

	serveCmd.Flags().Bool("debug", false, "Debug")

	return serveCmd
}

func doServe(cfg *internal.Config) error {
	h := &internal.Handler{
		Service:      &internal.Service{},
		PublicRouter: mux.NewRouter(),
	}

	h.SetRoutes()

	host := cfg.Public.Host
	if host == "" {
		host = "0.0.0.0"
	}

	port := cfg.Public.Port
	if port == 0 {
		port = int16(defaultPort)
	}

	addr := fmt.Sprintf("%v:%v", host, port)

	var router http.Handler = h.PublicRouter
	if cfg.CORS.Enable {
		opts := []handlers.CORSOption{}

		if cfg.CORS.AllowCredentials {
			opts = append(opts, handlers.AllowCredentials())
		}

		if cfg.CORS.MaxAge > 0 {
			opts = append(opts, handlers.MaxAge(cfg.CORS.MaxAge))
		}

		if len(cfg.CORS.AllowedHeaders) > 0 {
			opts = append(opts, handlers.AllowedHeaders(cfg.CORS.AllowedHeaders))
		}

		if len(cfg.CORS.AllowedOrigins) > 0 {
			opts = append(opts, handlers.AllowedOrigins(cfg.CORS.AllowedOrigins))
		}

		if len(cfg.CORS.AllowedMethods) > 0 {
			opts = append(opts, handlers.AllowedMethods(cfg.CORS.AllowedMethods))
		}

		if len(cfg.CORS.ExposedHeaders) > 0 {
			opts = append(opts, handlers.ExposedHeaders(cfg.CORS.ExposedHeaders))
		}

		router = handlers.CORS(opts...)(router)
	}

	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	tlsConfig := &tls.Config{}

	if r := cfg.Public.TLS.RequireClientCert; r {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if file := cfg.Public.TLS.CACert; file != "" {
		caCert, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.ClientCAs = caCertPool
		tlsConfig.BuildNameToCertificate()
	}

	if cfg.Debug {
		log.Println("Debug is enabled")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Starting server at %v", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}

	log.Println("server stopped")
	return nil
}
