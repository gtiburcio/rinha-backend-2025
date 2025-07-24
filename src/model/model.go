package model

import (
	"errors"
	"time"
)

type PaymentRequest struct {
	CorrelationID string    `json:"correlationId"`
	Amount        float64   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
}

func (p PaymentRequest) Valid() error {
	if p.CorrelationID == "" {
		return errors.New("invalid correlationId")
	}

	if p.Amount == float64(0) {
		return errors.New("invalid amount")
	}

	return nil
}

type PaymentSummaryDTO struct {
	Type   string
	Amount float64
}

type PaymentSummaryResponse struct {
	Default  PaymentSummaryResponseDetail `json:"default"`
	Fallback PaymentSummaryResponseDetail `json:"fallback"`
}

type PaymentSummaryResponseDetail struct {
	TotalRequests uint    `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

func BuildResponse(models []PaymentSummaryDTO) (p PaymentSummaryResponse) {
	def := PaymentSummaryResponseDetail{
		TotalRequests: 0,
		TotalAmount:   0,
	}
	fal := PaymentSummaryResponseDetail{
		TotalRequests: 0,
		TotalAmount:   0,
	}
	for _, m := range models {
		if m.Type == "default" {
			def.TotalAmount = def.TotalAmount + m.Amount
			def.TotalRequests = def.TotalRequests + 1
		} else {
			fal.TotalAmount = fal.TotalAmount + m.Amount
			fal.TotalRequests = fal.TotalRequests + 1
		}
	}

	p.Default = def
	p.Fallback = fal

	return
}
