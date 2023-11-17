package internal

import (
	"log"
	"net/http"
)

func (h *Handler) SetRoutes() {
	h.setDefaultHandlers()
	h.setAdminRoutes()
	h.setPublicRoutes()
	h.setHealthRoutes()
}

func (h *Handler) setDefaultHandlers() {
	h.PublicRouter.MethodNotAllowedHandler = http.HandlerFunc(h.MethodNotAllowed)
	h.PublicRouter.NotFoundHandler = http.HandlerFunc(h.NotFound)
}

const (
	HeaderTransactionID = "X-Transaction-Id"
)

func (h *Handler) setPublicRoutes() {
	r := h.PublicRouter.NewRoute().Subrouter()

	r.HandleFunc("/v1/transactions", h.CreateTransaction).Methods(http.MethodPost)
	r.HandleFunc("/v1/transactions", h.ListTransactions).Methods(http.MethodGet)
	r.HandleFunc("/v1/transactions/{id}", h.GetTransaction).Methods(http.MethodGet)
	r.HandleFunc("/v1/transactions/{id}", h.GetTransaction).Methods(http.MethodPut)
	r.HandleFunc("/v1/transactions/{id}/confirm", h.GetTransaction).Methods(http.MethodPost)
	r.HandleFunc("/v1/receiving-methods/validate", h.ValidateReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1/receiving-methods/retrieve", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1/balances", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1/products", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1/services", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1/operators", h.RetrieveReceivingMethod).Methods(http.MethodPost)

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Public routes")
			w.Header().Add(HeaderTransactionID, "1")

			h.ServeHTTP(w, r)
		})
	})
}

func (h *Handler) setAdminRoutes() {
	r := h.PublicRouter.NewRoute().Subrouter()

	r.HandleFunc("/v1/providers", h.CreateTransaction).Methods(http.MethodPost).Name("CreateProvider")
	r.HandleFunc("/v1/providers", h.ListTransactions).Methods(http.MethodGet).Name("ListProviders")
	r.HandleFunc("/v1/providers/{id}", h.ListTransactions).Methods(http.MethodDelete)
	r.HandleFunc("/v1/providers/{id}", h.ListTransactions).Methods(http.MethodPut)

	r.HandleFunc("/v1/services", h.CreateTransaction).Methods(http.MethodPost).Name("CreateService")
	r.HandleFunc("/v1/services", h.ListTransactions).Methods(http.MethodGet).Name("ListServices")
	r.HandleFunc("/v1/services/{id}", h.ListTransactions).Methods(http.MethodDelete)
	r.HandleFunc("/v1/services/{id}", h.ListTransactions).Methods(http.MethodPut)

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Admin routes")

			h.ServeHTTP(w, r)
		})
	})
}

func (h *Handler) setHealthRoutes() {
	h.PublicRouter.HandleFunc("/v1/ping", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/healthz/ready", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/healthz/live", h.ListTransactions).Methods(http.MethodGet)
}
