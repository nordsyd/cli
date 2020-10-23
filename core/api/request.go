package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

var apiURL = "http://localhost:5000"

// Get sends an HTTP GET request to the Nordsyd API
func Get(url string) (string, error) {
	//response, error := http.Get(apiURL + url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL+url, nil)

	if viper.Get("JWT") != "" {
		req.Header.Set("Authorization", "Bearer "+viper.Get("JWT").(string))
	}

	response, error := client.Do(req)

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

	// requestBody := strings.NewReader(string(requestString))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL+url, bytes.NewBuffer(requestString))

	if viper.Get("JWT") != "" {
		req.Header.Set("Authorization", "Bearer "+viper.Get("JWT").(string))
	}

	req.Header.Set("Content-Type", "application/json")

	response, error := client.Do(req)

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
