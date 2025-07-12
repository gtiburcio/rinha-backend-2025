package job

import (
	"context"
	"errors"
	"log"
	"rinha-backend-2025-gtiburcio/src/apperror"
	"rinha-backend-2025-gtiburcio/src/model"
	"time"
)

type (
	Job struct {
		queue      chan jobDTO
		maxRetries uint
		retryDelay time.Duration
		useCase    Usecase
	}

	Usecase interface {
		ProcessPayment(ctx context.Context, pr model.PaymentRequest) error
	}

	jobDTO struct {
		numRetries uint
		data       model.PaymentRequest
	}
)

func NewJob(useCase Usecase) Job {
	return Job{
		queue:      make(chan jobDTO, 20000),
		maxRetries: 5,
		retryDelay: time.Second * 2,
		useCase:    useCase,
	}
}

func (j Job) Enqueue(p model.PaymentRequest) error {
	dto := jobDTO{
		numRetries: 0,
		data:       p,
	}
	select {
	case j.queue <- dto:
		return nil
	default:
		return errors.New("fail to enqueue")
	}
}

func (j Job) Run() {
	for i := 0; i < 15; i++ {
		go j.Exec()
	}
}

func (j Job) Exec() {
	for dto := range j.queue {
		ctx := context.WithoutCancel(context.Background())
		err := j.useCase.ProcessPayment(ctx, dto.data)
		if err != nil && !apperror.IsIgnorableError(err) {
			j.enqueueToRetry(dto)
		}
	}
}

func (j Job) enqueueToRetry(dto jobDTO) {
	if j.maxRetries > dto.numRetries {
		dto.numRetries++

		log.Default().Printf("enqueue to retry payment with id %s numRetry %d", dto.data.CorrelationID, dto.numRetries)

		time.Sleep(j.retryDelay)

		j.queue <- dto
	} else {
		log.Default().Printf("discard payment with id %s numRetries %d", dto.data.CorrelationID, dto.numRetries)
	}
}
