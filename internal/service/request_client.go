package service

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RequestClient struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
}

// NewClient клиент для запросов
func NewClient(rl *rate.Limiter) *RequestClient {
	return &RequestClient{
		client:      http.DefaultClient,
		Ratelimiter: rl,
	}
}

// NewClientWithLimitForOneSecond клиент с количеством запросов в секунду
func NewClientWithLimitForOneSecond(numRequest int) *RequestClient {
	limiter := rate.NewLimiter(rate.Every(1*time.Second), numRequest)
	return NewClient(limiter)
}

// Do отправка запроса
func (c *RequestClient) Do(req *http.Request) (*http.Response, error) {
	var err error
	ctx := context.Background()
	// блокирует вызов для соблюдения лимита
	err = c.Ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
