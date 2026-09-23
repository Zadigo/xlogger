package models

import (
	"time"
)

// File represents a log file with its metadata,
// including UUID, name, path, size, and last modified time.
type File struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64 `json:"size"`
	LastModified time.Time `json:"last_modified"`
}
