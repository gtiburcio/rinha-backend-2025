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
	Type        string
	Amount      float64
	RequestedAt time.Time
}

type PaymentSummaryResponse struct {
	Default  PaymentSummaryResponseDetail `json:"default"`
	Fallback PaymentSummaryResponseDetail `json:"fallback"`
}

type PaymentSummaryResponseDetail struct {
	TotalRequests uint    `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

func BuildResponse(models []PaymentSummaryDTO, from, to string) (p PaymentSummaryResponse, err error) {
	fromTime, err := time.Parse(time.RFC3339Nano, from)
	if err != nil {
		return p, err
	}

	toTime, err := time.Parse(time.RFC3339Nano, to)
	if err != nil {
		return p, err
	}

	def := PaymentSummaryResponseDetail{
		TotalRequests: 0,
		TotalAmount:   0,
	}
	fal := PaymentSummaryResponseDetail{
		TotalRequests: 0,
		TotalAmount:   0,
	}
	for _, m := range models {
		c1 := fromTime.Compare(m.RequestedAt)
		c2 := toTime.Compare(m.RequestedAt)
		if c1 == 1 || c2 == -1 {
			continue
		}

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
