package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func TestMain(m *testing.M) {

	h := os.Getenv("roboz_web_host")
	if h == "" {
		os.Setenv("roboz_web_host", "localhost")
	}
	p := os.Getenv("roboz_web_port")
	if p == "" {
		os.Setenv("roboz_web_port", "5000")
	}

	opts := godog.Options{
		Format:        "pretty", // or "progress" (less verbose) or "junit" (continuous integration)
		Paths:         []string{"features"},
		Output:        colors.Colored(os.Stdout),
		StopOnFailure: false,
		Strict:        true,
	}

	status := godog.TestSuite{
		Name:                "smoke",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	hf := healthFeature{httpClient: httpClient}
	ctx.Step(`^the controller asks roboz whether it is healthy$`, hf.theControllerAsksRobozWhetherItIsHealthy)
	ctx.Step(`^roboz replies ok$`, hf.robozRepliesOk)

	df := documentationFeature{httpClient: httpClient}
	ctx.Step(`^the controller asks to read the documentation$`, df.theControllerAsksToReadTheDocumentation)
	ctx.Step(`^the documentation page is displayed$`, df.theDocumentationPageIsDisplayed)

	ef := entryFeature{httpClient: httpClient}
	ctx.Step(`^the controller asks for a list of executions$`, ef.theControllerAsksForAListOfExecutions)
	ctx.Step(`^the number of commands is at least (\d+)$`, ef.theNumberOfCommandsIsAtLeast)
	ctx.Step(`^the controller sends (\d+) instructions to the robot$`, ef.theControllerSendsInstructionsToTheRobot)
	ctx.Step(`^a response is received which matches the instructions$`, ef.aResponseIsreceivedWhichMatchesTheInstructions)
	ctx.Step(`^the controller sends invalid commands to the robot$`, ef.theControllerSendsInvalidCommandsToTheRobot)
	ctx.Step(`^an error response is received$`, ef.anErrorResponseIsreceived)
}
