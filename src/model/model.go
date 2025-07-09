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
