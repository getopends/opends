package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/getopends/opends/pkg/handler"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	handler.Handler
	Config       *Config
	Service      *TransactionService
	PublicRouter *mux.Router

	DB *sqlx.DB
}

func (h *Handler) Ping(rw http.ResponseWriter, req *http.Request) {
	body, err := parseTransactionInput(req)
	if err != nil {
		h.SendJSON(rw, err)
		return
	}

	resp, err := h.Service.CreateTransaction(body)
	if err != nil {
		h.SendJSON(rw, err)
		return
	}

	h.SendJSON(rw, resp)
}

func (h *Handler) CreateTransaction(rw http.ResponseWriter, req *http.Request) {
	body, err := parseTransactionInput(req)
	if err != nil {
		h.SendJSON(rw, err)
		return
	}

	resp, err := h.Service.CreateTransaction(body)
	if err != nil {
		h.SendJSON(rw, err)
		return
	}

	h.SendJSON(rw, resp)
}

func (h *Handler) ListTransactions(rw http.ResponseWriter, req *http.Request) {
	opts, err := parseTransactionsOptions(req)
	if err != nil {
		h.SendJSON(rw, err)
		return
	}

	resp, err := h.Service.ListTransactions(opts)
	if err != nil {
		h.SendJSON(rw, err)
		return
	}

	h.SendJSON(rw, resp)
}

func (h *Handler) GetTransaction(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		h.SendJSON(rw, apiErr)
		return
	}

	resp, apiErr := h.Service.GetTransaction(id)
	if apiErr != nil {
		h.SendJSON(rw, apiErr)
		return
	}

	h.SendJSON(rw, resp)
}

func (h *Handler) ValidateReceivingMethod(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		h.SendJSON(rw, apiErr)
		return
	}

	resp, apiErr := h.Service.GetTransaction(id)
	if apiErr != nil {
		h.SendJSON(rw, apiErr)
		return
	}

	h.SendJSON(rw, resp)
}

func (h *Handler) RetrieveReceivingMethod(rw http.ResponseWriter, req *http.Request) {
	id, apiErr := parseTransactionID(req)
	if apiErr != nil {
		h.SendJSON(rw, apiErr)
		return
	}

	resp, apiErr := h.Service.GetTransaction(id)
	if apiErr != nil {
		h.SendJSON(rw, apiErr)
		return
	}

	h.SendJSON(rw, resp)
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
	vars := mux.Vars(req)

	id, ok := vars["id"]
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
	opts := &ListTransactionOptions{}

	if value := req.URL.Query().Get("external_id"); value != "" {
		opts.ExternalID = value
	}

	if value := req.URL.Query().Get("page"); value != "" {
		page, err := strconv.Atoi(value)
		if err != nil {
			return nil, &Problem{
				Detail: "hey",
			}
		}

		opts.Page = page
	}

	if value := req.URL.Query().Get("per_page"); value != "" {
		perPage, err := strconv.Atoi(value)
		if err != nil {
			return nil, &Problem{
				Detail: "hey",
			}
		}

		opts.PerPage = perPage
	}

	return opts, nil
}
