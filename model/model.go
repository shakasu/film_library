package model

type Film struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ReleaseDate string `json:"releaseDate"`
	Rating      int    `json:"rating"`
}

type Actor struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_Of_birth"`
}
