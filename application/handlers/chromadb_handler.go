package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
)

type IChromaHandler interface {
	Upsert(w http.ResponseWriter, r *http.Request)
}

type chromaHandler struct {
	ingest services.IIngestFile
}

func NewChromaHandler(ingest services.IIngestFile) IChromaHandler {
	return &chromaHandler{
		ingest: ingest,
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

	if err := c.ingest.IngestToChroma(ctx,
		requestBody.CollectionName,
		requestBody.ContentId,
		requestBody.Content,
		requestBody.Metadata,
		services.WithDefaultIngestOptions()); err != nil {
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
	w.WriteHeader(202)
	w.Write(b)
}
