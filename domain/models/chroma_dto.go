package models

type ChromaUpsertRequest struct {
	CollectionName string                 `json:"collection_name"`
	ContentId      string                 `json:"content_id"`
	Content        string                 `json:"content"`
	Metadata       map[string]interface{} `json:"metadata"`
}
