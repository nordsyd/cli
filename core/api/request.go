package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var apiURL = "https://api.nordsyd.dk"

// Get sends an HTTP GET request to the Nordsyd API
func Get(url string) (string, error) {
	response, error := http.Get(apiURL + url)

	if error != nil {
		return "", error
	}

	data, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	return string(data), nil
}

// Post sends an HTTP POST request to the Nordsyd API
func Post(url string, payload interface{}) (string, error) {
	requestString, _ := json.Marshal(payload)

	requestBody := strings.NewReader(string(requestString))

	response, error := http.Post(
		apiURL+url,
		"application/json; charset=UTF-8",
		requestBody,
	)

	if error != nil {
		return "", error
	}

	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		var responseMap map[string]interface{}

		json.Unmarshal(data, &responseMap)

		// Ensure there is an error field
		responseError, ok := responseMap["error"]

		if !ok {
			errorCode := GetErrorCode("NORDSYD_API_NOT_WORKING")

			return "", errors.New(errorCode.humanReadable)
		}

		errorCode := GetErrorCode(responseError.(string))

		return "", errors.New(errorCode.humanReadable)
	}

	response.Body.Close()

	return string(data), nil
}
