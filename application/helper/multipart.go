package helper

import (
	"mime/multipart"
	"net/http"
	"slices"
)

var (
	AllowArrayField = []string{}
)

type MultipartFile struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type MultipartRequest struct {
	Fields map[string]interface{}
	Files  map[string]*MultipartFile
}

const (
	MaxMultipartMemory = 5 << 20
)

func ParseMultipartRequest(r *http.Request) (*MultipartRequest, error) {
	if err := r.ParseMultipartForm(MaxMultipartMemory); err != nil {
		return nil, err
	}

	multipartRequest := &MultipartRequest{
		Fields: make(map[string]interface{}),
		Files:  make(map[string]*MultipartFile),
	}

	// extract text fields from form-data
	for key, val := range r.MultipartForm.Value {
		if len(val) > 0 {
			if slices.Contains(AllowArrayField, key) {
				multipartRequest.Fields[key] = val
			} else {
				multipartRequest.Fields[key] = val[0]
			}
		}
	}

	// extract fiels from form-data
	for key, val := range r.MultipartForm.File {
		if len(val) > 0 {
			file, err := val[0].Open()
			if err != nil {
				return nil, err
			}

			multipartRequest.Files[key] = &MultipartFile{
				File:       file,
				FileHeader: val[0],
			}
		}
	}

	return multipartRequest, nil
}
