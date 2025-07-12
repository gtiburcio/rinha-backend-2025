package repository

import (
	"context"
	"rinha-backend-2025-gtiburcio/src/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	summaryQuery = `select type, count(1) as "count", sum(amount) as "sum" from payment
			 where requested_at between $1 and $2
			 group by type;
	`
)

type (
	Repository struct {
		dbConn *pgxpool.Pool
	}
)

func NewRepository(dbConn *pgxpool.Pool) Repository {
	return Repository{
		dbConn: dbConn,
	}
}

func (r Repository) Save(ctx context.Context, pr model.PaymentRequest, strategyType string) error {
	_, err := r.dbConn.Exec(ctx, "insert into payment (correlation_id, amount, type, requested_at) values ($1, $2, $3, $4)",
		pr.CorrelationID, pr.Amount, strategyType, pr.RequestedAt,
	)

	return err
}

func (r Repository) FindSummary(ctx context.Context, from, to string) ([]model.PaymentSummaryDTO, error) {
	rows, err := r.dbConn.Query(ctx, summaryQuery, from, to)
	if err != nil {
		return nil, err
	}

	dtos := []model.PaymentSummaryDTO{}

	for rows.Next() {
		psm := model.PaymentSummaryDTO{}
		rows.Scan(&psm.Type, &psm.Count, &psm.Sum)
		dtos = append(dtos, psm)
	}

	return dtos, nil
}
