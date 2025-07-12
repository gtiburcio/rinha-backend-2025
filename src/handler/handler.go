package handler

import (
	"context"
	"io"
	"net/http"
	"rinha-backend-2025-gtiburcio/src/model"

	"github.com/goccy/go-json"
)

type (
	Handler struct {
		job     Job
		useCase UseCase
	}

	Job interface {
		Enqueue(p model.PaymentRequest) error
	}

	UseCase interface {
		FindSummary(ctx context.Context, from, to string) ([]model.PaymentSummaryDTO, error)
	}
)

func NewHandler(job Job, useCase UseCase) Handler {
	return Handler{
		job:     job,
		useCase: useCase,
	}
}

func (h Handler) HandleSavePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr := model.PaymentRequest{}
	if err := json.Unmarshal(body, &pr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := pr.Valid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.job.Enqueue(pr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) HandlePaymentSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	dtos, err := h.useCase.FindSummary(r.Context(), from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(model.BuildResponse(dtos))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}
