package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"igcinfo/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"igcinfo/api"
)

/*
Testing the first endpoint of the api, the index- handler.
*/
func TestIndexAPI(t *testing.T) {

	// Testing table containing slice of struct with all the diff scenerios
	TestTable := []struct {
		method        string
		url           string
		expStatusCode int
	}{
		{method: "GET", url: "https://igcinfoapi.herokuapp.com/api", expStatusCode: http.StatusOK},
		{method: "GET", url: "https://igcinfoapi.herokuapp.com", expStatusCode: http.StatusNotFound},
		{method: "GET", url: "https://igcinfoapi.herokuapp.com/apii", expStatusCode: http.StatusNotFound},
		{method: "POST", url: "https://igcinfoapi.herokuapp.com/api/", expStatusCode: http.StatusNotFound},
	}

	for idx, testCase := range TestTable {

		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Could't make get request %v", err)
		}

		writer := httptest.NewRecorder()
		api.Index(writer, req)

		response := writer.Result()

		if response.StatusCode != testCase.expStatusCode {
			t.Fatalf("Testcase(%v) Bad response code, expected %v, got %v", idx, testCase.expStatusCode, response.StatusCode)
		}
	}

}

/*
Test that ShowTrackInfo function works properly.
*/
func TestShowTrackInfo(t *testing.T) {

	TestTable := []struct {
		method        string
		url           string
		expStatusCode int
	}{
		{method: "GET", url: "localhost:8080/api/igc/1", expStatusCode: http.StatusNotFound},
		{method: "GET", url: "localhost:8080/api/igc/1/pilot/rr", expStatusCode: http.StatusNotFound},
		{method: "GET", url: "localhost:8080/api/igc/rr", expStatusCode: http.StatusNotFound},
		{method: "POST", url: "localhost:8080/api/igc/1", expStatusCode: http.StatusNotFound},
	}

	for idx, testCase := range TestTable {

		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Could't make get request %v", err)
		}

		writer := httptest.NewRecorder()
		api.Index(writer, req)

		response := writer.Result()

		if response.StatusCode != testCase.expStatusCode {
			t.Fatalf("Testcase(%v) Bad response code, expected %v, got %v", idx, testCase.expStatusCode, response.StatusCode)
		}
	}
}

/*
Testing RegAndShowTrackIds
*/
func TestRegAndShowTrackIds(t *testing.T) {
	// Testing table containing slice of struct with all the diff scenarios
	TestTable := []struct {
		method        string
		url           string
		expStatusCode int
		body          map[string]string
	}{
		{method: "GET", url: "localhost:8080/api/igc", expStatusCode: http.StatusOK},
		{method: "POST", url: "localhost:8080/api/igc", expStatusCode: http.StatusOK},
	}

	for idx, testCase := range TestTable {

		if testCase.method == "GET" {
			req, err := http.NewRequest(testCase.method, testCase.url, nil)
			if err != nil {
				t.Errorf("Could't make request, %v", err)
			}

			writer := httptest.NewRecorder()
			api.RegAndShowTrackIds(writer, req)

			response := writer.Result()

			if response.StatusCode != testCase.expStatusCode {
				t.Errorf("Testcase(%v) Bad response code, expected %v, got %v", idx, testCase.expStatusCode, response.StatusCode)
			}
		}

		if testCase.method == "POST" {

			testCase.body = make(map[string]string)
			testCase.body["url"] = "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
			b, err := json.Marshal(&testCase.body)
			buf := bytes.NewBuffer(b)

			if err != nil {
				t.Errorf("Error marshalling json, err: %v", err)
			}

			req, err := http.NewRequest(testCase.method, testCase.url, buf)
			if err != nil {
				t.Errorf("Could't make request, %v", err)
			}

			writer := httptest.NewRecorder()
			api.RegAndShowTrackIds(writer, req)

			response := writer.Result()

			if response.StatusCode != testCase.expStatusCode {
				t.Errorf("Testcase(%v) Bad response code, expected %v, got %v", idx, testCase.expStatusCode, response.StatusCode)
			}
		}
	}
}

/*
Testing GetUptime
*/
func TestGetUptime(t *testing.T) {
	startTimes := make(map[int]time.Time)
	startTimes[1] = time.Date(2018, 9, 30, 12, 30, 40, 2, time.UTC)
	startTimes[2] = time.Date(2018, 10, 12, 9, 0, 40, 2, time.UTC)
	startTimes[3] = time.Date(2018, 10, 10, 17, 1, 2, 0, time.UTC)
	startTimes[4] = time.Date(2017, 9, 0, 0, 50, 12, 0, time.UTC)

	// Outputting so i can verify that the time format is correct
	for idx, times := range startTimes {
		res := model.GetUptime(times)
		fmt.Println(idx, res)
	}

}
