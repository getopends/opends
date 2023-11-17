package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/getopends/opends/internal/handler"
	"github.com/gorilla/mux"
)

type Handler struct {
	handler.Base
	Service      *Service
	PublicRouter *mux.Router
}

func (h *Handler) SetRoutes() {
	h.setDefaultHandlers()
	h.setAdminRoutes()
	h.setPublicRoutes()
	h.setProbes()
}

func (h *Handler) setDefaultHandlers() {
	h.PublicRouter.MethodNotAllowedHandler = http.HandlerFunc(h.MethodNotAllowed)
	h.PublicRouter.NotFoundHandler = http.HandlerFunc(h.NotFound)
}

func (h *Handler) setPublicRoutes() {
	h.PublicRouter.HandleFunc("/v1alpha1/transactions", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/transactions", h.ListTransactions).Methods(http.MethodGet)
	h.PublicRouter.HandleFunc("/v1alpha1/transactions/{id}", h.GetTransaction).Methods(http.MethodGet)
	h.PublicRouter.HandleFunc("/v1alpha1/receiving-methods/validate", h.ValidateReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/receiving-methods/retrieve", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/balances", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/products", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/services", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/operators", h.RetrieveReceivingMethod).Methods(http.MethodPost)
}

func (h *Handler) setAdminRoutes() {
	h.PublicRouter.HandleFunc("/v1alpha1/transactions", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/transactions", h.ListTransactions).Methods(http.MethodGet)
	h.PublicRouter.HandleFunc("/v1alpha1/transactions/{id}", h.GetTransaction).Methods(http.MethodGet)
	h.PublicRouter.HandleFunc("/v1alpha1/receiving-methods/validate", h.ValidateReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/receiving-methods/retrieve", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/balances", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/products", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/services", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/v1alpha1/operators", h.RetrieveReceivingMethod).Methods(http.MethodPost)
}

func (h *Handler) setProbes() {
	h.PublicRouter.HandleFunc("/healthz/ready", h.CreateTransaction).Methods(http.MethodPost)
	h.PublicRouter.HandleFunc("/healthz/live", h.ListTransactions).Methods(http.MethodGet)
}

func (h *Handler) CreateTransaction(rw http.ResponseWriter, req *http.Request) {
	body, err := parseTransactionInput(req)
	if err != nil {
		h.JSON(rw, err)
		return
	}

	resp, err := h.Service.CreateTransaction(body)
	if err != nil {
		h.JSON(rw, err)
		return
	}

	h.JSON(rw, resp)
}

func (h *Handler) ListTransactions(rw http.ResponseWriter, req *http.Request) {
	opts, err := parseTransactionsOptions(req)
	if err != nil {
		h.JSON(rw, err)
		return
	}

	resp, err := h.Service.ListTransactions(opts)
	if err != nil {
		h.JSON(rw, err)
		return
	}

	h.JSON(rw, resp)
}

func (h *Handler) GetTransaction(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		h.JSON(rw, apiErr)
		return
	}

	resp, apiErr := h.Service.GetTransaction(id)
	if apiErr != nil {
		h.JSON(rw, apiErr)
		return
	}

	h.JSON(rw, resp)
}

func (h *Handler) ValidateReceivingMethod(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		h.JSON(rw, apiErr)
		return
	}

	resp, apiErr := h.Service.GetTransaction(id)
	if apiErr != nil {
		h.JSON(rw, apiErr)
		return
	}

	h.JSON(rw, resp)
}

func (h *Handler) RetrieveReceivingMethod(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		h.JSON(rw, apiErr)
		return
	}

	resp, apiErr := h.Service.GetTransaction(id)
	if apiErr != nil {
		h.JSON(rw, apiErr)
		return
	}

	h.JSON(rw, resp)
}

func (h Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/problem+json")

	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(&Problem{
		Detail: "method not allowed",
	})
}

func (h Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/problem+json")

	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(&Problem{
		Detail: "method not allowed",
	})
}

func parseTransactionInput(req *http.Request) (*CreateTransactionInput, *Problem) {
	body := &CreateTransactionInput{}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, &Problem{Detail: "hey"}
	}

	if err := json.Unmarshal(data, body); err != nil {
		return nil, &Problem{Detail: "hey"}
	}

	return body, nil
}

func parseTransactionID(req *http.Request) (uint64, *Problem) {
	v := mux.Vars(req)

	id, ok := v["id"]
	if !ok {
		return 0, &Problem{Detail: "hey"}
	}

	newId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, &Problem{Detail: "hey"}
	}

	return newId, nil
}

func parseTransactionsOptions(req *http.Request) (*ListTransactionOptions, *Problem) {
	externalID := req.URL.Query().Get("external_id")
	if externalID == "" {
		return nil, &Problem{
			Detail: "hey",
		}
	}

	return &ListTransactionOptions{
		ExternalID: externalID,
	}, nil
}
