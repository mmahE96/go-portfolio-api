package models

// User schema of the user table
type Article struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	Date        string `json:"date"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

type Project struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	YearDone    string `json:"year_done"`
	ClientName  string `json:"client_name"`
	Category    string `json:"category"`
	Date        string `json:"date"`
}
