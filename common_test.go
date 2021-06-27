package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func videoShouldBeDownloaded(videoId string) error {
	if err := getVideo(videoId); err != nil {
		return err
	}

	if videoResponse["status"] != "DOWNLOADED" {
		return errors.New("invalid status. got: " + videoResponse["status"].(string))
	}

	return nil
}

func getVideo(videoId string) error {
	resp, err := http.Get("http://" + appUrl + "/" + videoId)

	if err != nil {
		return errors.New("error while trying to get video. err: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid status code received after sending video get request. got: " + strconv.Itoa(resp.StatusCode))
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&videoResponse)
	}

	return nil
}

func getCut(videoId string, cutStart string, cutEnd string) error {
	if err := getVideo(videoId); err != nil {
		return err
	}

	cuts, ok := videoResponse["cuts"].([]interface{})

	if !ok {
		return errors.New("cant parse video cuts")
	}

	var foundCut map[string]interface{}
	for _, cut := range cuts {
		if cutAsserted, ok := cut.(map[string]interface{}); ok {
			if cutAsserted["start"] == cutStart && cutAsserted["end"] == cutEnd {
				foundCut = cutAsserted
				break
			}
		}
	}

	if foundCut == nil {
		return errors.New("cannot find cut")
	}

	cutResponse = foundCut

	return nil
}

func waitMinutes(minutes int) error {
	<-time.After(time.Duration(rand.Intn(minutes)) * time.Minute)
	return nil
}
