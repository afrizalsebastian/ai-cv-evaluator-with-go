package helper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrBodyRequestNil       = errors.New("request body is empty")
	ErrReadBodyRequest      = errors.New("error when read body request")
	ErrUnmarshalBodyRequest = errors.New("error when unmarshal body request")
)

func ParseJSONBodyRequest[T any](r *http.Request) (*T, error) {
	var data T

	if r.Body == nil {
		return nil, ErrBodyRequestNil
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, ErrReadBodyRequest
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, ErrUnmarshalBodyRequest
	}

	return &data, nil
}
