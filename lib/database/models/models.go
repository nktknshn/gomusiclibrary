package models

import "time"

type TagID int64

type AudioTag struct {
	ID     TagID
	FileID FileID

	Name  string
	Value string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
