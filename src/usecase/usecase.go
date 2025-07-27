package usecase

import (
	"context"
	"rinha-backend-2025-gtiburcio/src/model"
	"time"
)

type (
	UseCase struct {
		client         Client
		repository     Repository
		concurrentdata Concurrentdata
	}

	Client interface {
		SavePayment(ctx context.Context, pr model.PaymentRequest) (string, error)
	}

	Repository interface {
		Save(ctx context.Context, pr model.PaymentRequest, strategyType string) error

		FindSummary(ctx context.Context, from, to string) ([]model.PaymentSummaryDTO, error)

		FindAll() ([]model.PaymentSummaryDTO, error)
	}

	Concurrentdata interface {
		GetData() []model.PaymentSummaryDTO
	}
)

func NewUseCase(client Client, repository Repository, concurrentdata Concurrentdata) UseCase {
	return UseCase{
		client:         client,
		repository:     repository,
		concurrentdata: concurrentdata,
	}
}

func (u UseCase) ProcessPayment(ctx context.Context, pr model.PaymentRequest) error {
	pr.RequestedAt = time.Now()
	t, err := u.client.SavePayment(ctx, pr)
	if err != nil {
		return err
	}

	if err = u.repository.Save(ctx, pr, t); err != nil {
		return err
	}

	return nil
}

func (u UseCase) FindSummary() []model.PaymentSummaryDTO {
	return u.concurrentdata.GetData()
}

func (u UseCase) FindAll() ([]model.PaymentSummaryDTO, error) {
	return u.repository.FindAll()
}
