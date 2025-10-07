package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
	chromaclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/chroma-client"
)

type IChromaHandler interface {
	Upsert(w http.ResponseWriter, r *http.Request)
}

type chromaHandler struct {
	chroma chromaclient.IChromaClient
}

func NewChromaHandler(chroma chromaclient.IChromaClient) IChromaHandler {
	return &chromaHandler{
		chroma: chroma,
	}
}

func (c *chromaHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var requestBody models.ChromaUpsertRequest
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &requestBody); err != nil {
		fmt.Printf("error when parse request body: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusBadRequest, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJson)
		return
	}

	if err := c.chroma.Upsert(ctx, requestBody.CollectionName, requestBody.ContentId, requestBody.Content, requestBody.Metadata); err != nil {
		fmt.Printf("error when upsert to chroma: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusInternalServerError, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respJson)
		return
	}

	resp := models.CreateWebResponse("Success", http.StatusOK, nil)
	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
