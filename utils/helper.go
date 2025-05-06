package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/sirupsen/logrus"
)

type HTTPClientService struct {
	AuthToken string
	Method    string
	URL       string
}

var ExponentialBackoffRetry = func(maxRetries int, operation func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		logrus.Errorf("Attempt %d: Operation failed: %s", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
		}
	}
	return err
}

func (hcp *HTTPClientService) prepareRequest(body interface{}) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		logrus.Errorf("Failed to marshal body in MakeAPICall: %v", err)
		return nil, err
	}

	req, err := http.NewRequest(hcp.Method, hcp.URL, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if hcp.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", hcp.AuthToken))
	}

	return req, nil
}

func (hcp *HTTPClientService) readResponseBody(resp *http.Response, result interface{}) error {
	err := json.NewDecoder(resp.Body).Decode(result)

	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}
	return nil
}

func (hcp *HTTPClientService) MakeAPICall(body interface{}, result interface{}) error {
	req, err := hcp.prepareRequest(body)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: config.AppConfig.TIMEOUT,
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Failed to make request: %v", err)
		return err
	}

	defer resp.Body.Close()
	return hcp.readResponseBody(resp, result)
}
