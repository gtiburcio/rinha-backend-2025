package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"rinha-backend-2025-gtiburcio/src/apperror"
	"rinha-backend-2025-gtiburcio/src/model"
	"time"

	"github.com/goccy/go-json"
)

type Client struct {
	client             http.Client
	defaultServiceURL  string
	fallbackServiceURL string
}

func NewClient() Client {
	return Client{
		client: http.Client{
			Timeout: time.Second * 2,
		},
		defaultServiceURL:  getBaseURL(os.Getenv("DEFAULT_HOST")),
		fallbackServiceURL: getBaseURL(os.Getenv("FALLBACK_HOST")),
	}
}

func (c Client) SavePayment(ctx context.Context, pr model.PaymentRequest) (string, error) {
	j, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}

	err = c.execCall(c.defaultServiceURL, j)
	if err == nil {
		return "default", nil
	}

	if apperror.IsIgnorableError(err) {
		return "", err
	}

	err = c.execCall(c.fallbackServiceURL, j)
	if err == nil {
		return "fallback", nil
	}

	return "", err
}

func (c Client) execCall(baseURL string, payload []byte) error {
	url := baseURL + "/payments"
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error to call payments api: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnprocessableEntity {
			return apperror.NewAppError(true, "unprocessable")
		}
		return fmt.Errorf("error to call payments api: %v", err)
	}

	return nil
}

func getBaseURL(host string) string {
	return fmt.Sprintf("http://%s", host)
}
