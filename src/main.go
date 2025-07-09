package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rinha-backend-2025-gtiburcio/src/config"
	"rinha-backend-2025-gtiburcio/src/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	dbConn *pgxpool.Pool
}

func main() {
	app := App{
		dbConn: config.NewDBConfig().DBConn,
	}

	http.HandleFunc("/payments", app.handleSavePayment)

	fmt.Println("Server starting on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func (a App) handleSavePayment(w http.ResponseWriter, r *http.Request) {
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

	if err := a.savePayment(r.Context(), pr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a App) savePayment(ctx context.Context, pr model.PaymentRequest) error {
	t, err := a.execClientCall(ctx, &pr, false)
	if err != nil {
		return err
	}

	_, err = a.dbConn.Exec(ctx, "insert into payment (correlation_id, amount, type, requested_at) values ($1, $2, $3, $4)", pr.CorrelationID, pr.Amount, t, pr.RequestedAt)
	if err != nil {
		return err
	}

	return nil
}

func (a App) execClientCall(ctx context.Context, pr *model.PaymentRequest, fallback bool) (string, error) {
	pr.RequestedAt = time.Now().UTC()

	client := http.Client{
		Timeout: time.Second * 2,
	}

	j, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}

	url := getBaseURL(fallback) + "/payments"
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(j))
	if err != nil || resp.StatusCode != http.StatusOK {
		if !fallback {
			return a.execClientCall(ctx, pr, true)
		}
		return "", fmt.Errorf("error to call payments api: %v", err)
	}

	t := "default"

	if fallback {
		t = "fallback"
	}

	return t, nil
}

func getBaseURL(fallback bool) string {
	port := 8001
	if fallback {
		port = 8002
	}
	return fmt.Sprintf("http://localhost:%d", port)
}
