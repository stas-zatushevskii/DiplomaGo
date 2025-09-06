package order

import (
	"encoding/json"
	"fmt"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"io"
	"net/http"
	"time"
)

func (o *ServiceOrder) SendRequest(orderNumber string) (*models.AccrualResponse, error) {
	var response models.AccrualResponse
	URL := fmt.Sprintf("%s/api/orders/%s", o.config.Accrual.Address, orderNumber)
	o.logger.Warn(URL)
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 4,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type RequestSender interface {
	SendRequest(orderNumber string) (*models.AccrualResponse, error)
}

// RequestWithRetry retrying send request, increase timeout += 2, max retry == 3
func (o *ServiceOrder) RequestWithRetry(sender RequestSender, data string) (*models.AccrualResponse, error) {
	retryCount := 3
	timeout := 1

	response, err := sender.SendRequest(data)
	if err != nil {
		return nil, err
	}

	if response.Status == constants.OrderStatusProcessed {
		return response, nil
	}

	if !isRetryable(response.Status) {
		return nil, fmt.Errorf("response status is %s", response.Status)
	}

	for i := 0; i < retryCount; i++ {
		o.logger.Info(fmt.Sprintf("Retryable error: %v. Retrying in %d seconds...\n", response.Status, timeout))
		time.Sleep(time.Duration(timeout) * time.Second)
		timeout += 2

		response, err = sender.SendRequest(data)
		if err != nil {
			return nil, err
		}

		if response.Status == constants.OrderStatusProcessed {
			return response, nil
		}
		if !isRetryable(response.Status) {
			return nil, fmt.Errorf("response status is %s", response.Status)
		}
	}

	return nil, fmt.Errorf("retried %d times, last error: %s", retryCount, response.Status)
}

// isRetryable Check if response form request is retryable
func isRetryable(status constants.OrderStatus) bool {
	if status == constants.OrderStatusProcessing || status == constants.OrderStatusNew {
		return true
	}
	return false
}
