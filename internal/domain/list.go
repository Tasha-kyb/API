package domain

import "time"

type List struct {
	ID			string		`json:"id"`
	Title 		string		`json:"title"`
	CreatedAt	time.Time	`json:"created_at"`
}