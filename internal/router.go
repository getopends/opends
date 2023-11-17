package internal

import "net/http"

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

func (h *Handler) setPublicRoutes() {
	h.PublicRouter.HandleFunc("/v1/transactions", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/transactions", h.ListTransactions).Methods(http.MethodGet)
	h.PublicRouter.HandleFunc("/v1alpha1/transactions/{id}", h.GetTransaction).Methods(http.MethodGet)
	h.PublicRouter.HandleFunc("/v1/receiving-methods/validate", h.ValidateReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/receiving-methods/retrieve", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/balances", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/products", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/services", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/operators", h.RetrieveReceivingMethod).Methods(http.MethodPost)
}

func (h *Handler) setAdminRoutes() {
	h.PublicRouter.HandleFunc("/v1/transactions", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/transactions", h.ListTransactions).Methods(http.MethodGet)
	h.PublicRouter.HandleFunc("/v1/transactions/{id}", h.GetTransaction).Methods(http.MethodGet)
	h.PublicRouter.HandleFunc("/v1/receiving-methods/validate", h.ValidateReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/receiving-methods/retrieve", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/balances", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/products", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/services", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1/operators", h.RetrieveReceivingMethod).Methods(http.MethodPost)
}

func (h *Handler) setHealthRoutes() {
	h.PublicRouter.HandleFunc("/v1/ping", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/healthz/ready", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/healthz/live", h.ListTransactions).Methods(http.MethodGet)
}
