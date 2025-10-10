package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/api"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/helper"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
)

type IUploadDocumentController interface {
	PostDocument(ctx context.Context, r *http.Request) api.WebResponse
}

type uploadDocumentController struct {
	uploadDocumentService services.IUploadDocumentService
}

func NewUploadDocumenController(uploadDocumentService services.IUploadDocumentService) IUploadDocumentController {
	return &uploadDocumentController{
		uploadDocumentService: uploadDocumentService,
	}
}

func (u *uploadDocumentController) PostDocument(ctx context.Context, r *http.Request) api.WebResponse {
	multipartRequest, err := helper.ParseMultipartRequest(r)
	if err != nil {
		log.Println("error when parse body request")
		return api.CreateWebResponse("Invalid request", http.StatusBadRequest, nil, nil)
	}

	request := helper.MultipartToUploadDocumentRequest(multipartRequest)

	if structErr := helper.ValidateParams(ctx, request); structErr != nil {
		log.Println("error validation on body request")
		return api.CreateWebResponse("Invalid request", http.StatusBadRequest, nil, structErr)
	}

	resp := u.uploadDocumentService.SaveUploadedDocument(ctx, request)
	return resp

}
