package chromaclient

import (
	"context"
	"fmt"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/domain/models"
	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/amikos-tech/chroma-go/pkg/embeddings"
)

type ChromaNotFoundRecord struct {
	Query          string
	CollectionName string
}

func (c *ChromaNotFoundRecord) Error() string {
	return fmt.Sprintf("Not Found record at %s with query %s", c.CollectionName, c.Query)
}

type IChromaClient interface {
	Upsert(ctx context.Context, collectionName, id, content string, metadata map[string]interface{}) error
	Query(ctx context.Context, collectionName, query string, topK int) ([]models.ChromaSearchResult, error)
}

type chromaClient struct {
	cli chroma.Client
}

func NewChromaClient(ctx context.Context, chromaUrl string) (IChromaClient, error) {
	client, err := chroma.NewHTTPClient(
		chroma.WithBaseURL(chromaUrl),
		chroma.WithDatabaseAndTenant(chroma.DefaultDatabase, chroma.DefaultTenant),
	)

	if err != nil {
		return nil, err
	}

	return &chromaClient{cli: client}, nil
}

func (c *chromaClient) Upsert(ctx context.Context, collectionName, id, content string, metadata map[string]interface{}) error {
	collection, err := c.cli.GetOrCreateCollection(ctx, collectionName, chroma.WithEmbeddingFunctionCreate(embeddings.NewConsistentHashEmbeddingFunction()))
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	var metaAttributes []*chroma.MetaAttribute
	for k, v := range metadata {
		switch val := v.(type) {
		case string:
			metaAttributes = append(metaAttributes, chroma.NewStringAttribute(k, val))
		case int:
			metaAttributes = append(metaAttributes, chroma.NewIntAttribute(k, int64(val)))
		case float64:
			metaAttributes = append(metaAttributes, chroma.NewFloatAttribute(k, val))
		default:
			fmt.Printf("Warning: Unsupported metadata type for key '%s'\n", k)
		}
	}

	if err := collection.Upsert(ctx,
		chroma.WithIDs(chroma.DocumentID(id)),
		chroma.WithTexts(content),
		chroma.WithMetadatas(chroma.NewDocumentMetadata(metaAttributes...))); err != nil {
		return fmt.Errorf("failed to upsert document: %w", err)
	}

	return nil
}

func (c *chromaClient) Query(ctx context.Context, collectionName, query string, topK int) ([]models.ChromaSearchResult, error) {
	embedding := embeddings.NewConsistentHashEmbeddingFunction()
	collection, err := c.cli.GetCollection(ctx, collectionName, chroma.WithEmbeddingFunctionGet(embedding))
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	embeddingQuery, _ := embedding.EmbedQuery(ctx, query)
	resp, err := collection.Query(ctx,
		chroma.WithNResults(topK),
		chroma.WithQueryEmbeddings(embeddingQuery),
		chroma.WithIncludeQuery(chroma.IncludeDocuments, chroma.IncludeMetadatas),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %w", err)
	}

	if resp == nil || len(resp.GetIDGroups()) == 0 || len(resp.GetIDGroups()[0]) == 0 {
		return nil, &ChromaNotFoundRecord{CollectionName: collectionName, Query: query}
	}

	var results []models.ChromaSearchResult
	idGroup := resp.GetIDGroups()[0]
	docsGroup := resp.GetDocumentsGroups()[0]

	for i, id := range idGroup {
		result := models.ChromaSearchResult{
			Id:   string(id),
			Text: docsGroup[i].ContentString(),
		}

		results = append(results, result)
	}

	return results, nil
}
