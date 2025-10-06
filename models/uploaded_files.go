package models

import "time"

type UploadedFile struct {
	ID       string
	Path     string
	Filename string
	Uploaded time.Time
}

type UploadedFileResponse struct {
	FileId string `json:"file_id"`
}
