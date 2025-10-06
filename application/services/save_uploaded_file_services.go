package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrSaveFile     = errors.New("file error to save")
	ErrDeleteFile   = errors.New("file failed to delete")
)

type ILocalFileStore interface {
	SaveUploadedFile(fileId string, file multipart.File, header *multipart.FileHeader) (*models.UploadedFile, error)
	Get(id string) (*models.UploadedFile, error)
	Delete(id string) error
}

type localFileStore struct {
	basePath string
	meta     map[string]*models.UploadedFile
	mu       sync.RWMutex
}

func NewLocalFileStore(base string) (ILocalFileStore, error) {
	_, err := os.Stat(base)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(base, 0o755); err != nil {
			fmt.Printf("Warning: Failed to create base directory %s: %v\n", base, err)
			return nil, err
		}
	}

	return &localFileStore{
		basePath: base,
		meta:     make(map[string]*models.UploadedFile),
	}, nil
}

func (s *localFileStore) SaveUploadedFile(fileId string, file multipart.File, header *multipart.FileHeader) (*models.UploadedFile, error) {
	filename := header.Filename

	safeFilename := filepath.Clean(filename)
	if safeFilename == "." || safeFilename == "/" {
		safeFilename = filename
	}

	_, err := os.Stat(filepath.Join(s.basePath, fileId))
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(s.basePath, fileId), 0o755); err != nil {
			fmt.Printf("Warning: Failed to create directory %s: %v\n", fileId, err)
			return nil, err
		}
	}

	dest := filepath.Join(s.basePath, fileId, safeFilename)

	out, err := os.Create(dest)
	if err != nil {
		return nil, ErrSaveFile
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		_ = os.Remove(dest)
		return nil, ErrSaveFile
	}

	if err := out.Sync(); err != nil {
		return nil, ErrSaveFile
	}

	uf := &models.UploadedFile{
		ID:       fileId,
		Path:     dest,
		Filename: filename,
		Uploaded: time.Now(),
	}

	s.mu.Lock()
	s.meta[fileId+safeFilename] = uf
	s.mu.Unlock()

	return uf, nil
}

func (s *localFileStore) Get(fileId string) (*models.UploadedFile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.meta[fileId]
	if !ok {
		return nil, ErrFileNotFound
	}
	return v, nil
}

func (s *localFileStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	uf, ok := s.meta[id]
	if !ok {
		return ErrFileNotFound
	}

	if err := os.Remove(uf.Path); err != nil {
		if !os.IsNotExist(err) {
			return ErrDeleteFile
		}
	}

	delete(s.meta, id)
	return nil
}
