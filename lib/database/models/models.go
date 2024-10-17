package models

import "time"

type FileID int64
type TagID int64

type File struct {
	ID         FileID
	Path       string
	Size       int64
	Sha256Hash string

	Ctime int64
	Mtime int64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type FileSlice []File

type AudioTag struct {
	ID     TagID
	FileID FileID

	Name  string
	Value string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
