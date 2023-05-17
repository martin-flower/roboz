package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type entryFeature struct {
	httpClient    *http.Client
	response      *http.Response
	entryResponse entryResponse
	listResponses []listResponse
	sent          int
}

type listResponse struct {
	Commands int `json:"commands"`
}

type entryResponse struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Commands  int       `json:"commands"` // note specification is for commmands, not commands
	Result    int       `json:"result"`
	Duration  float64   `json:"duration"`
}

type Coordinate struct {
	X int
	Y int
}

type Command struct {
	Direction string
	Steps     int
}

type enterPost struct {
	Start    Coordinate
	Commands []Command
}

// --- when ---

func (ef *entryFeature) theControllerAsksForAListOfExecutions() (err error) {
	var response *http.Response
	response, err = http.Get(fmt.Sprintf("http://%s:%s/list", os.Getenv("roboz_web_host"), os.Getenv("roboz_web_port")))
	if err != nil {
		return
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&ef.listResponses)
	return
}

func (ef *entryFeature) theControllerSendsInstructionsToTheRobot(sent int) (err error) {

	tosend := enterPost{}
	tosend.Start = Coordinate{X: 4, Y: 5}
	for i := 0; i < sent; i++ {
		tosend.Commands = append(tosend.Commands, Command{Direction: "north", Steps: 2})
	}

	tosendJSON, err := json.Marshal(tosend)

	if err != nil {
		return fmt.Errorf("failed to marshall json %w", err)
	}

	var request *http.Request
	request, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/developer-test/enter-path", os.Getenv("roboz_web_host"), os.Getenv("roboz_web_port")), bytes.NewBuffer(tosendJSON))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("failed to create request %w", err)
	}
	var response *http.Response
	response, err = ef.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to post request %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code %d", response.StatusCode)
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&ef.entryResponse)
	if err != nil {
		return fmt.Errorf("failed to marshall response %w", err)
	}

	ef.sent = sent
	return
}

func (ef *entryFeature) theControllerSendsInvalidCommandsToTheRobot() error {

	tosend := enterPost{}
	tosend.Start = Coordinate{X: 4, Y: 5}
	tosend.Commands = []Command{Command{Direction: "kkk", Steps: -4}, Command{Direction: "mmm", Steps: 0}}

	tosendJSON, err := json.Marshal(tosend)

	if err != nil {
		return fmt.Errorf("failed to marshall json %w", err)
	}

	var request *http.Request
	request, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/developer-test/enter-path", os.Getenv("roboz_web_host"), os.Getenv("roboz_web_port")), bytes.NewBuffer(tosendJSON))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("failed to create request %w", err)
	}
	var response *http.Response
	response, err = ef.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to post request %w", err)
	}

	ef.response = response

	return nil
}

// --- then --- (asserts response)

func (ef *entryFeature) aResponseIsreceivedWhichMatchesTheInstructions() error {
	if ef.sent != ef.entryResponse.Commands {
		return fmt.Errorf("unexpected %d commands, received %d", ef.sent, ef.entryResponse.Commands)
	}
	return nil
}

func (ef *entryFeature) anErrorResponseIsreceived() error {
	if ef.response.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("unexpected status code %d", ef.response.StatusCode)
	}
	return nil
}

func (ef *entryFeature) theNumberOfCommandsIsAtLeast(expected int) error {
	if ef.sent > ef.entryResponse.Commands {
		return fmt.Errorf("expected at least %d commands, received %d", ef.sent, ef.entryResponse.Commands)
	}
	return nil
}
