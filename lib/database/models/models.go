package models

import "time"

type AudioTagID int64

type AudioTag struct {
	ID     AudioTagID
	FileID FileID

	Name  string
	Value string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type AudioTagSlice []AudioTag
