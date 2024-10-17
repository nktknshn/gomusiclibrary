package models

import "time"

type FileID int64

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

func (fs FileSlice) IDs() []FileID {
	return MapFileSlice(fs, func(f File) FileID {
		return f.ID
	})
}

func MapFileSlice[T any](fs FileSlice, fn func(File) T) []T {
	res := make([]T, len(fs))
	for idx, f := range fs {
		res[idx] = fn(f)
	}
	return res
}
