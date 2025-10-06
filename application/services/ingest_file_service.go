package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	chromaclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/chroma-client"
	"github.com/ledongthuc/pdf"
)

var (
	ErrOpenPdfFile               = errors.New("error open pdf file")
	ErrIngestCollectionNameEmpty = errors.New("error ingest collection name empty")
	ErrIngestDocIdEmpty          = errors.New("error ingest docId empty")
	ErrIngestContentEmpty        = errors.New("error ingest content empty")
	ErrIngestNoContentGenerated  = errors.New("error ingest no content generated")
)

type ChunkingConfig struct {
	WordsPerChunk int
	OverlapWords  int
}

func WithDefaultChunkConfig() ChunkingConfig {
	return ChunkingConfig{
		WordsPerChunk: 150,
		OverlapWords:  20,
	}
}

type IngestOptions struct {
	ChunkingConfig ChunkingConfig
	BatchSize      int
}

func WithDefaultIngestOptions() IngestOptions {
	return IngestOptions{
		ChunkingConfig: WithDefaultChunkConfig(),
		BatchSize:      10,
	}
}

// Ingest
type IIngestFile interface {
	ExtractTextFromPdf(path string) (string, error)
	ChunkText(text string, config ChunkingConfig) []string
	IngestToChroma(ctx context.Context, collectionName, docId, content string, metadata map[string]interface{}, option IngestOptions) error
}

type ingestFile struct {
	chroma chromaclient.IChromaClient
}

func NewIngestFile(chroma chromaclient.IChromaClient) IIngestFile {
	return &ingestFile{
		chroma: chroma,
	}
}

func (i *ingestFile) ExtractTextFromPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", ErrOpenPdfFile
	}
	defer f.Close()

	var sb strings.Builder
	totalPage := r.NumPage()

	for page := 1; page <= totalPage; page++ {
		p := r.Page(page)
		if p.V.IsNull() {
			continue
		}

		content, err := p.GetPlainText(nil)
		if err != nil {
			fmt.Println("error when get the text from file")
			continue
		}

		sb.WriteString(content)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (i *ingestFile) ChunkText(text string, config ChunkingConfig) []string {
	// sanitize config
	if config.WordsPerChunk <= 0 {
		config.WordsPerChunk = 150
	}
	if config.OverlapWords < 0 {
		config.OverlapWords = 0
	}
	if config.OverlapWords >= config.WordsPerChunk {
		config.OverlapWords = config.WordsPerChunk / 2
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var chunks []string
	step := config.WordsPerChunk - config.OverlapWords

	for i := 0; i < len(words); i += step {
		end := i + config.WordsPerChunk
		if end > len(words) {
			end = len(words)
		}

		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, chunk)

		if end == len(words) {
			break
		}
	}

	return chunks
}

func (i *ingestFile) IngestToChroma(ctx context.Context, collectionName, docId, content string, metadata map[string]interface{}, options IngestOptions) error {
	if collectionName == "" {
		return ErrIngestCollectionNameEmpty
	}
	if docId == "" {
		return ErrIngestDocIdEmpty
	}
	if content == "" {
		return ErrIngestContentEmpty
	}

	chunks := i.ChunkText(content, options.ChunkingConfig)
	if len(chunks) == 0 {
		return errors.New("no chunks generated from content")
	}

	// Add chunk metadata
	for idx, chunk := range chunks {
		recID := fmt.Sprintf("%s_chunk_%d", docId, idx)

		chunkMetadata := make(map[string]interface{})
		for k, v := range metadata {
			chunkMetadata[k] = v
		}
		chunkMetadata["chunk_index"] = strconv.Itoa(idx)
		chunkMetadata["total_chunks"] = strconv.Itoa(len(chunks))
		chunkMetadata["document_id"] = docId

		if err := i.chroma.Upsert(ctx, collectionName, recID, chunk, chunkMetadata); err != nil {
			return errors.New("failed to upsert chunk")
		}
	}

	return nil
}
