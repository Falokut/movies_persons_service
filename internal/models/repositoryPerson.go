package models

import "time"

type RepositoryPerson struct {
	ID         string    `db:"id" json:"id"`
	FullnameRU string    `db:"fullname_ru" json:"fullname_ru"`
	FullnameEN string    `db:"fullname_en" json:"fullname_en,omitempty"`
	Birthday   time.Time `db:"birthday" json:"birthday,omitempty"`
	Sex        string    `db:"sex" json:"sex,omitempty"`
	PhotoID    string    `db:"photo_id" json:"photo_url,omitempty"`
}
