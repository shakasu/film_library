package model

type Film struct {
	Id          int
	Name        string
	Description string
	ReleaseDate string
	rating      int
}

type Actor struct {
	Id          int
	Name        string
	gender      string
	dateOfBirth string
}
