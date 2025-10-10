package models

import "mime/multipart"

type UploadDocumentRequest struct {
	CvFile           multipart.File        `form:"cv_file" validate:"required"`
	CvFileHeader     *multipart.FileHeader `form:"cv_file_header"`
	ReportFile       multipart.File        `form:"report_file" validate:"required"`
	ReportFileHeader *multipart.FileHeader `form:"report_file_header"`
}

type UploadedDocumentResponse struct {
	FileId string `json:"file_id"`
}
