package api

import (
	"encoding/json"
)

// GetUserSites gets all sites associated with the JWT
func GetUserSites() []Site {
	response, _ := Get("/site")

	var sites []Site

	json.Unmarshal([]byte(response), &sites)

	return sites
}
