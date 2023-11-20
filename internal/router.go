package internal

import (
	"log"
	"net/http"
)

func (h *Handler) RegisterRoutes() {
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
	r := h.PublicRouter.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/transactions", h.CreateTransaction).Methods(http.MethodPost)
	r.HandleFunc("/transactions", h.ListTransactions).Methods(http.MethodGet)
	r.HandleFunc("/transactions/{id}", h.GetTransaction).Methods(http.MethodGet)
	r.HandleFunc("/transactions/{id}", h.GetTransaction).Methods(http.MethodPut)
	r.HandleFunc("/transactions/{id}/confirm", h.GetTransaction).Methods(http.MethodPost)
	r.HandleFunc("/receiving-methods/validate", h.ValidateReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/receiving-methods/retrieve", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/balances", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/products", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/services", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/operators", h.RetrieveReceivingMethod).Methods(http.MethodPost)

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Public routes")
			w.Header().Add(HeaderTransactionID, "1")

			h.ServeHTTP(w, r)
		})
	})
}

func (h *Handler) setAdminRoutes() {
	r := h.PublicRouter.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/providers", h.CreateTransaction).Methods(http.MethodPost).Name("CreateProvider")
	r.HandleFunc("/providers", h.ListTransactions).Methods(http.MethodGet).Name("ListProviders")
	r.HandleFunc("/providers/{id}", h.ListTransactions).Methods(http.MethodDelete)
	r.HandleFunc("/providers/{id}", h.ListTransactions).Methods(http.MethodPut)

	r.HandleFunc("/services", h.CreateTransaction).Methods(http.MethodPost).Name("CreateService")
	r.HandleFunc("/services", h.ListTransactions).Methods(http.MethodGet).Name("ListServices")
	r.HandleFunc("/services/{id}", h.ListTransactions).Methods(http.MethodDelete)
	r.HandleFunc("/services/{id}", h.ListTransactions).Methods(http.MethodPut)

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
