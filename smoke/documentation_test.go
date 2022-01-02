package main

import (
	"fmt"
	"net/http"
	"os"
)

type documentationFeature struct {
	httpClient *http.Client
	statusCode int
}

func (df *documentationFeature) theControllerAsksToReadTheDocumentation() (err error) {

	var request *http.Request
	request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s", os.Getenv("roboz_web_host"), os.Getenv("roboz_web_port")), nil)
	if err != nil {
		return fmt.Errorf("failed to create request %w", err)
	}
	var response *http.Response
	response, err = df.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to get request %w", err)
	}

	df.statusCode = response.StatusCode

	return
}

func (df *documentationFeature) theDocumentationPageIsDisplayed() error {
	if df.statusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", df.statusCode)
	}
	return nil
}
