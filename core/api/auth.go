package api

import (
	"encoding/json"
)

// GetJWT gets a JWT token given an email and password
func GetJWT(email string, password string) (string, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	response, error := Post("/token/issue", payload)

	if error != nil {
		return "", error
	}

	var responseMap map[string]interface{}

	json.Unmarshal([]byte(response), &responseMap)

	return responseMap["jwt"].(string), nil
}
