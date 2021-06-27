package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

var (
	cutPayload  = map[string]interface{}{}
	cutResponse = map[string]interface{}{}
)

func iSendACutRequestWithStartAndEnd(arg1, arg2 string) error {
	cutPayload["start"] = arg1
	cutPayload["end"] = arg2

	cutBytes, err := json.Marshal([]interface{}{cutPayload})

	if err != nil {
		return errors.New("cant serialize cut request. err: " + err.Error())
	}

	req, err := http.NewRequest(http.MethodPut, "http://"+appUrl+"/videos/"+videoResponse["id"].(string)+"/queue", bytes.NewBuffer(cutBytes))

	if err != nil {
		return errors.New("cant create cut request. err: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.New("cant send cut request. err: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("invalid status code received after sending cut request. got: " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

func iVerifyIfTheCutWasReallyMade() error {
	if err := getCut(videoResponse["id"].(string), cutPayload["start"].(string), cutPayload["end"].(string)); err != nil {
		return err
	}

	if cutResponse["status"] != "CUT" {
		return errors.New("invalid status. got: " + cutResponse["status"].(string))
	}

	return nil
}

func iWaitSecondsSoTheCutRequestCanBeProcessedAndCut(arg1 int) error {
	return waitSeconds(arg1)
}

func theCutShouldBeInAProcessingState() error {
	if err := getCut(videoResponse["id"].(string), cutPayload["start"].(string), cutPayload["end"].(string)); err != nil {
		return err
	}

	if !(cutResponse["status"] == "QUEUE" || cutResponse["status"] == "CUTTING" || cutResponse["status"] == "RENDERING" || cutResponse["status"] == "CUT") {
		return errors.New("invalid status. got: " + cutResponse["status"].(string))
	}

	return nil
}

func theVideoIsDownloadedAndReadyToCut(arg1 string) error {
	return videoShouldBeDownloaded(arg1)
}

func InitializeCutScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I send a cut request with start "([^"]*)" and end "([^"]*)"$`, iSendACutRequestWithStartAndEnd)
	ctx.Step(`^I verify if the cut was really made$`, iVerifyIfTheCutWasReallyMade)
	ctx.Step(`^I wait (\d+) minutes so the cut request can be processed and cut$`, iWaitSecondsSoTheCutRequestCanBeProcessedAndCut)
	ctx.Step(`^the cut should be in a processing state$`, theCutShouldBeInAProcessingState)
	ctx.Step(`^the video "([^"]*)" is downloaded and ready to cut$`, theVideoIsDownloadedAndReadyToCut)
}
