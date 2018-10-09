package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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
		{method: "GET", url: "https://igcinfoapi.herokuapp.com/api/", expStatusCode: http.StatusOK},
		{method: "GET", url: "https://igcinfoapi.herokuapp.com", expStatusCode: http.StatusNotFound},
		{method: "GET", url: "https://igcinfoapi.herokuapp.com/apii", expStatusCode: http.StatusNotFound},
		{method: "POST", url: "https://igcinfoapi.herokuapp.com/api/", expStatusCode: http.StatusMethodNotAllowed},
	}

	for idx, testCase := range TestTable {

		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Could't make get request %v", err)
		}

		writer := httptest.NewRecorder()
		Index(writer, req)

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
		Index(writer, req)

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
			RegAndShowTrackIds(writer, req)

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

			//os.Stdout.Write(b)
			if err != nil {
				fmt.Errorf("Error marshalling json, err: %v", err)
			}

			req, err := http.NewRequest(testCase.method, testCase.url, buf)
			if err != nil {
				t.Errorf("Could't make request, %v", err)
			}

			writer := httptest.NewRecorder()
			RegAndShowTrackIds(writer, req)

			response := writer.Result()

			if response.StatusCode != testCase.expStatusCode {
				t.Errorf("Testcase(%v) Bad response code, expected %v, got %v", idx, testCase.expStatusCode, response.StatusCode)
			}
		}
	}
}








