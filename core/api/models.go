package api

// Site is the struct for the site model
type Site struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	DateCreated string `json:"date_created"`
}
