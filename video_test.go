package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/cucumber/godog"
)

var (
	videoPayload  = map[string]interface{}{}
	videoResponse = map[string]interface{}{}
)

func iSendTheDownloadRequestWithTheOutputPathAndTheVideoItselfAsOpeningAndEnding(arg1 string) error {
	videoResponse = map[string]interface{}{}
	videoPayload["opts"] = map[string]interface{}{
		"output_dir":   arg1,
		"opening_path": filepath.Join(arg1, videoPayload["id"].(string), videoPayload["id"].(string)+".mp4"),
		"ending_path":  filepath.Join(arg1, videoPayload["id"].(string), videoPayload["id"].(string)+".mp4"),
	}

	videoBytes, err := json.Marshal(videoPayload)

	if err != nil {
		return errors.New("cant serialize video request. err: " + err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+appUrl+"/queue", bytes.NewBuffer(videoBytes))

	if err != nil {
		return errors.New("cant create video request. err: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.New("cant send video request. err: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("invalid status code received after sending download request. got: " + strconv.Itoa(resp.StatusCode))
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&videoResponse)
	}

	return nil
}

func iVerifyIfTheVideoWasReallyDownloaded() error {
	return videoShouldBeDownloaded(videoResponse["id"].(string))
}

func iWaitSecondsSoTheVideoCanBeProcessedAndDownloaded(arg1 int) error {
	return waitSeconds(arg1)
}

func iWantToDownloadTheVideoThatContainsTheId(arg1 string) error {
	videoPayload["id"] = arg1
	return nil
}

func theVideoShouldBeInAProcessingState() error {
	if err := getVideo(videoResponse["id"].(string)); err != nil {
		return err
	}

	if !(videoResponse["status"] == "QUEUE" || videoResponse["status"] == "DOWNLOADING" || videoResponse["status"] == "DOWNLOADED") {
		return errors.New("invalid status. got: " + videoResponse["status"].(string))
	}

	return nil
}

func InitializeVideoScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I send the download request with the output path "([^"]*)" and the video itself as opening and ending$`, iSendTheDownloadRequestWithTheOutputPathAndTheVideoItselfAsOpeningAndEnding)
	ctx.Step(`^I verify if the video was really downloaded$`, iVerifyIfTheVideoWasReallyDownloaded)
	ctx.Step(`^I wait (\d+) minutes so the video can be processed and downloaded$`, iWaitSecondsSoTheVideoCanBeProcessedAndDownloaded)
	ctx.Step(`^I want to download the video that contains the id "([^"]*)"$`, iWantToDownloadTheVideoThatContainsTheId)
	ctx.Step(`^the video should be in a processing state$`, theVideoShouldBeInAProcessingState)
}
