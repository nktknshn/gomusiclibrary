package models

import (
	"time"

	"github.com/nktknshn/gomusiclibrary/lib/util/colutil"
)

type FileID int64

type File struct {
	ID         FileID
	Path       string
	Size       int64
	Sha256Hash string

	Mtime time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type FileSlice []File

func (fs FileSlice) IDs() []FileID {
	return colutil.MapSlice(fs, func(f File) FileID {
		return f.ID
	})
}

func (fs FileSlice) Paths() []string {
	return colutil.MapSlice(fs, func(f File) string {
		return f.Path
	})

}

type FileMapID map[FileID]File
type FileMapPath map[string]File
