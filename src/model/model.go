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
	Type  string
	Count uint
	Sum   float64
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
	for _, m := range models {
		detail := PaymentSummaryResponseDetail{
			TotalRequests: m.Count,
			TotalAmount:   m.Sum,
		}
		if m.Type == "default" {
			p.Default = detail
		} else {
			p.Fallback = detail
		}
	}

	return
}
