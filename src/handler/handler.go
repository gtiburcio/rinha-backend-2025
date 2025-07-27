package handler

import (
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
		FindSummary() []model.PaymentSummaryDTO

		FindAll() ([]model.PaymentSummaryDTO, error)
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

	pr := model.PaymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
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

	dtos, _ := h.useCase.FindAll()

	responseList, err := model.BuildResponse(dtos, from, to)
	if err != nil {
		http.Error(w, "invalid dates", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(responseList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}
