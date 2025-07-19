package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"rinha-backend-2025-gtiburcio/src/apperror"
	"rinha-backend-2025-gtiburcio/src/model"
	"strconv"
	"time"

	"github.com/goccy/go-json"
)

type Client struct {
	defaultClient      http.Client
	fallbackClient     http.Client
	retryTimes         int
	retryDelay         time.Duration
	defaultServiceURL  string
	fallbackServiceURL string
}

func NewClient() Client {
	retryTimes, err := strconv.ParseInt(os.Getenv("RETRY_TIMES"), 10, 64)
	if err != nil {
		retryTimes = 5
	}

	retryDelay, err := strconv.ParseInt(os.Getenv("RETRY_DELAY"), 10, 64)
	if err != nil {
		retryDelay = 500
	}

	defaultTimeout, err := strconv.ParseInt(os.Getenv("DEFAULT_TIMEOUT"), 10, 64)
	if err != nil {
		defaultTimeout = 200
	}

	fallbackTimeout, err := strconv.ParseInt(os.Getenv("FALLBACK_TIMEOUT"), 10, 64)
	if err != nil {
		fallbackTimeout = 500
	}

	return Client{
		defaultClient: http.Client{
			Timeout: time.Millisecond * time.Duration(defaultTimeout),
		},
		fallbackClient: http.Client{
			Timeout: time.Millisecond * time.Duration(fallbackTimeout),
		},
		retryTimes:         int(retryTimes),
		retryDelay:         time.Millisecond * time.Duration(retryDelay),
		defaultServiceURL:  getBaseURL(os.Getenv("DEFAULT_HOST")),
		fallbackServiceURL: getBaseURL(os.Getenv("FALLBACK_HOST")),
	}
}

func (c Client) SavePayment(ctx context.Context, pr model.PaymentRequest) (string, error) {
	j, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}

	for i := 0; i < c.retryTimes; i++ {
		err = c.execCall(c.defaultServiceURL, j, false)
		if err == nil {
			return "default", nil
		}

		if apperror.IsIgnorableError(err) {
			return "", err
		}

		time.Sleep(c.retryDelay)
	}

	err = c.execCall(c.fallbackServiceURL, j, true)
	if err == nil {
		return "fallback", nil
	}

	return "", err
}

func (c Client) execCall(baseURL string, payload []byte, isFallback bool) error {
	client := c.defaultClient
	if isFallback {
		client = c.fallbackClient
	}
	resp, err := client.Post(baseURL, "application/json", bytes.NewBuffer(payload))
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
	return fmt.Sprintf("http://%s/payments", host)
}
