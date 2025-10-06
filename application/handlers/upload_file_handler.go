package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
	"github.com/google/uuid"
)

const (
	MAX_FILE_SIZE = 5 << 20 // 5 MB
)

type IUploadFileHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

type uploadFileHandler struct {
	fileStore services.ILocalFileStore
}

func NewUploadHandler(filestore services.ILocalFileStore) IUploadFileHandler {
	return &uploadFileHandler{
		fileStore: filestore,
	}
}

func (u *uploadFileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(MAX_FILE_SIZE); err != nil {
		fmt.Printf("error when parse form-data: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusBadRequest, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJson)
		return
	}

	directoryId := uuid.New().String()

	cvFile, cvHeader, err := r.FormFile("cv_file")
	if err != nil {
		fmt.Printf("error when get cv file: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusBadRequest, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJson)
		return
	}
	defer cvFile.Close()

	cvHeader.Filename = "cv_file.pdf"
	_, err = u.fileStore.SaveUploadedFile(directoryId, cvFile, cvHeader)
	if err != nil {
		fmt.Printf("error when save cv file: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusInternalServerError, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respJson)
		return
	}

	reportFile, reportHeader, err := r.FormFile("report_file")
	if err != nil {
		fmt.Printf("error when get report file: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusBadRequest, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJson)
		return
	}
	defer reportFile.Close()

	reportHeader.Filename = "report_file.pdf"
	_, err = u.fileStore.SaveUploadedFile(directoryId, reportFile, reportHeader)
	if err != nil {
		fmt.Printf("error when save report file: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusInternalServerError, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respJson)
		return
	}

	uploadResponse := &models.UploadedFileResponse{
		FileId: directoryId,
	}
	resp := models.CreateWebResponse("Success", http.StatusOK, uploadResponse)
	respJson, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}
