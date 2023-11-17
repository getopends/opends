package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	srv *Service
}

func (h *Handler) CreateTransaction(rw http.ResponseWriter, req *http.Request) {
	body, err := parseCreateTransactionInput(req)
	if err != nil {
		sendJSON(rw, err)
		return
	}

	resp, err := h.srv.CreateTransaction(body)
	if err != nil {
		sendJSON(rw, err)
		return
	}

	sendJSON(rw, resp)
}

func (h *Handler) ListTransactions(rw http.ResponseWriter, req *http.Request) {
	opts, err := parseListTransactionsOptions(req)
	if err != nil {
		sendJSON(rw, err)
		return
	}

	resp, err := h.srv.ListTransactions(opts)
	if err != nil {
		sendJSON(rw, err)
		return
	}

	sendJSON(rw, resp)
}

func (h *Handler) GetTransaction(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		sendJSON(rw, apiErr)
		return
	}

	resp, apiErr := h.srv.GetTransaction(id)
	if apiErr != nil {
		sendJSON(rw, apiErr)
		return
	}

	sendJSON(rw, resp)
}

func (h *Handler) ValidateReceivingMethod(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		sendJSON(rw, apiErr)
		return
	}

	resp, apiErr := h.srv.GetTransaction(id)
	if apiErr != nil {
		sendJSON(rw, apiErr)
		return
	}

	sendJSON(rw, resp)
}

func (h *Handler) RetrieveReceivingMethod(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		sendJSON(rw, apiErr)
		return
	}

	resp, apiErr := h.srv.GetTransaction(id)
	if apiErr != nil {
		sendJSON(rw, apiErr)
		return
	}

	sendJSON(rw, resp)
}

func parseCreateTransactionInput(req *http.Request) (*CreateTransactionInput, *Problem) {
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

func parseListTransactionsOptions(req *http.Request) (*ListTransactionOptions, *Problem) {
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

func sendJSON(rw http.ResponseWriter, a any) error {
	return json.NewEncoder(rw).Encode(a)
}

type Problem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int64  `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

type Service struct{}

type ServiceOptions struct{}

func (s Service) CreateTransaction(body *CreateTransactionInput) (*Transaction, *Problem) {
	return &Transaction{
		ID:         1,
		ExternalID: body.ExternalID,
	}, nil
}

func (s Service) GetTransaction(id uint64) (*Transaction, *Problem) {
	return &Transaction{
		ID: id,
	}, nil
}

type ListTransactionOptions struct {
	ExternalID string `json:"external_id"`
}

type ListProductsOptions struct {
	OperatorID string `json:"operator_id"`
	ServiceID  string `json:"service_id"`
}

func (s Service) ListTransactions(opts *ListTransactionOptions) ([]Transaction, *Problem) {
	return []Transaction{
		{
			ID:         1,
			ExternalID: "1",
		},
	}, nil
}

type Transaction struct {
	ID         uint64 `json:"id"`
	ExternalID string `json:"external_id"`
}

type CreateTransactionInput struct {
	ExternalID string `json:"external_id"`
}

type ValidateReceivingMethodInput struct {
	SendingCustomer   *SendingCustomer   `json:"sending_customer,omitempty"`
	SendingBusiness   *SendingBusiness   `json:"sending_business,omitempty"`
	ReceivingCustomer *ReceivingCustomer `json:"receiving_customer,omitempty"`
	ReceivingBusiness *ReceivingBusiness `json:"receiving_business,omitempty"`
	ReceivingMethod   *ReceivingMethod   `json:"receiving_method,omitempty"`
}

type RetrieveReceivingMethodInput struct {
	ReceivingMethod *ReceivingMethod `json:"receiving_method,omitempty"`
}

type SendingBusiness struct {
	Name string `json:"name"`
}

type SendingCustomer struct {
	FirstName string `json:"first_name"`
}

type ReceivingBusiness struct {
	Name string `json:"name"`
}

type ReceivingCustomer struct {
	Firstname string `json:"firstname"`
}

type ReceivingMethod struct {
	CardNumber  string `json:"card_number"`
	PhoneNumber string `json:"phone_number"`
}

func main() {
	h := &Handler{
		srv: &Service{},
	}

	r := mux.NewRouter()

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/problem+json")

		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(&Problem{
			Detail: "method not allowed",
		})
	})

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/problem+json")

		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(&Problem{
			Detail: "not found",
		})
	})

	r.HandleFunc("/v1alpha1/transactions", h.CreateTransaction).Methods(http.MethodPost)
	r.HandleFunc("/v1alpha1/transactions", h.ListTransactions).Methods(http.MethodGet)
	r.HandleFunc("/v1alpha1/transactions/{id}", h.GetTransaction).Methods(http.MethodGet)
	r.HandleFunc("/v1alpha1/receiving-methods/validate", h.ValidateReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1alpha1/receiving-methods/retrieve", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1alpha1/balances", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1alpha1/products", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1alpha1/services", h.RetrieveReceivingMethod).Methods(http.MethodPost)
	r.HandleFunc("/v1alpha1/operators", h.RetrieveReceivingMethod).Methods(http.MethodPost)

	srv := http.Server{
		Addr:    ":12345",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

type Response struct {
	Header http.Header
	Body   io.Reader
}
