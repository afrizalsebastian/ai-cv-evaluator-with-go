package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/api"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/domain/models"
	"github.com/google/uuid"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrSaveFile     = errors.New("file error to save")
	ErrDeleteFile   = errors.New("file failed to delete")
)

type IUploadDocumentService interface {
	SaveUploadedDocument(context.Context, *models.UploadDocumentRequest) api.WebResponse
}

type uploadDocumentService struct {
	basePath string
}

func NewUploadDocumentService(basePath string) IUploadDocumentService {
	return &uploadDocumentService{
		basePath: basePath,
	}
}

func (u *uploadDocumentService) SaveUploadedDocument(ctx context.Context, req *models.UploadDocumentRequest) api.WebResponse {
	folderId := uuid.New().String()
	errChan := make(chan error, 2)
	// save cv file

	go func() {
		fmt.Println("upload cv")
		req.CvFileHeader.Filename = "cv_file.pdf"
		if err := u.saveToPath(folderId, req.CvFile, req.CvFileHeader); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	go func() {
		fmt.Println("upload report")
		req.ReportFileHeader.Filename = "report_file.pdf"
		if err := u.saveToPath(folderId, req.ReportFile, req.ReportFileHeader); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	fmt.Println("upload done")

	err := <-errChan
	if err != nil {
		log.Println("error when save document")
		return api.CreateWebResponse("Error when save user document", http.StatusInternalServerError, nil, nil)
	}

	resp := &models.UploadedDocumentResponse{
		FileId: folderId,
	}

	return api.CreateWebResponse("Success", http.StatusOK, resp, nil)

}

func (u *uploadDocumentService) saveToPath(folderId string, file multipart.File, header *multipart.FileHeader) error {
	filename := header.Filename

	safeFilename := filepath.Clean(filename)
	if safeFilename == "." || safeFilename == "/" {
		safeFilename = filename
	}

	_, err := os.Stat(filepath.Join(u.basePath, folderId))
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(u.basePath, folderId), 0o755); err != nil {
			log.Printf("Warning: Failed to create directory %s: %v\n", folderId, err)
			return err
		}
	}

	dest := filepath.Join(u.basePath, folderId, safeFilename)

	out, err := os.Create(dest)
	if err != nil {
		return ErrSaveFile
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		_ = os.Remove(dest)
		return ErrSaveFile
	}

	if err := out.Sync(); err != nil {
		return ErrSaveFile
	}

	return nil
}
