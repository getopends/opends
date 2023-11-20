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
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

var (
	defaultPort = int16(13000)
	defaultHost = "0.0.0.0"
)

func serveCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:  "serve",
		RunE: runCmd(doServe),
	}

	serveCmd.Flags().Bool("debug", false, "Debug")

	return serveCmd
}

func configureCORS(router http.Handler, cfg *internal.CORSOptions) http.Handler {
	if cfg.Enable {
		opts := []handlers.CORSOption{}

		if cfg.AllowCredentials {
			opts = append(opts, handlers.AllowCredentials())
		}

		if cfg.MaxAge > 0 {
			opts = append(opts, handlers.MaxAge(cfg.MaxAge))
		}

		if len(cfg.AllowedHeaders) > 0 {
			opts = append(opts, handlers.AllowedHeaders(cfg.AllowedHeaders))
		}

		if len(cfg.AllowedOrigins) > 0 {
			opts = append(opts, handlers.AllowedOrigins(cfg.AllowedOrigins))
		}

		if len(cfg.AllowedMethods) > 0 {
			opts = append(opts, handlers.AllowedMethods(cfg.AllowedMethods))
		}

		if len(cfg.ExposedHeaders) > 0 {
			opts = append(opts, handlers.ExposedHeaders(cfg.ExposedHeaders))
		}

		router = handlers.CORS(opts...)(router)
	}

	return router
}

func configureDB(cfg *internal.DatabaseOptions) (*sqlx.DB, error) {
	driver := cfg.Driver
	if driver == "" {
		driver = "postgres"
	}

	db, err := sqlx.Connect(driver, cfg.DSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func configureTLS(srv *http.Server, cfg *internal.ServerOptions) error {
	tlsConfig := &tls.Config{}

	if r := cfg.TLS.RequireClientCert; r {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if file := cfg.TLS.CACert; file != "" {
		caCert, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.ClientCAs = caCertPool
		tlsConfig.BuildNameToCertificate()
	}

	srv.TLSConfig = tlsConfig

	return nil
}

func newServer(r *mux.Router, cfg *internal.ServerOptions, cors *internal.CORSOptions) (*http.Server, error) {
	host := cfg.Host
	if host == "" {
		host = defaultHost
	}

	port := cfg.Port
	if port == 0 {
		port = defaultPort
	}

	addr := fmt.Sprintf("%v:%v", host, port)

	router := configureCORS(r, cors)

	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := configureTLS(&srv, cfg); err != nil {
		return nil, err
	}

	return &srv, nil
}

func doServe(ctx context.Context, _ *cobra.Command, _ []string, cfg *internal.Config) error {
	db, err := configureDB(&cfg.Database)
	if err != nil {
		return err
	}

	h := &internal.Handler{
		Service:      &internal.TransactionService{},
		PublicRouter: mux.NewRouter(),
		DB:           db,
	}

	h.RegisterRoutes()

	srv, err := newServer(h.PublicRouter, &cfg.Public, &cfg.CORS)
	if err != nil {
		return err
	}

	defer srv.Close()
	defer db.Conn(ctx)

	if cfg.Debug {
		log.Println("Debug is enabled")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Starting server at %v", srv.Addr)
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
