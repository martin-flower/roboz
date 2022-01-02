package main

import (
	"fmt"
	"net/http"
	"os"
)

type healthFeature struct {
	httpClient *http.Client
	statusCode int
}

func (hf *healthFeature) theControllerAsksRobozWhetherItIsHealthy() (err error) {

	var request *http.Request
	request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/health", os.Getenv("roboz_web_host"), os.Getenv("roboz_web_port")), nil)
	if err != nil {
		return fmt.Errorf("failed to create request %w", err)
	}
	var response *http.Response
	response, err = hf.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to get request %w", err)
	}
	if err != nil {
		return
	}
	hf.statusCode = response.StatusCode
	return
}

func (hf *healthFeature) robozRepliesOk() error {
	if hf.statusCode != http.StatusOK {
		return fmt.Errorf("expected status code %d, received %d", http.StatusOK, hf.statusCode)
	}
	return nil
}
