package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var (
	httpClient = &http.Client{
		Timeout: time.Second * 5,
	}
	appUrl string = "localhost:8080"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	InitializeVideoScenario(ctx)
	InitializeCutScenario(ctx)
}
func TestMain(m *testing.M) {

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  []string{"features/video.feature", "features/cut.feature"},
		},
	}.Run()

	os.Exit(status)
}
