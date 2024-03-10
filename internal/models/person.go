package models

import "time"

type Person struct {
	ID         string
	FullnameRU string
	FullnameEN string
	Birthday   time.Time
	Sex        string
	PhotoURL   string
}
