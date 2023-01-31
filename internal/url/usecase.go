package url

import (
	"context"

	"github.com/pkg/errors"
)

type (
	// ShortenUseCase shortens a http url. It implements Shortener interface.
	ShortenUseCase struct {
		generator  CodeGenerator
		repository Persister
	}
)

// NewShortenURLUseCase returns an instance ShortenUseCase.
func NewShortenURLUseCase(generator CodeGenerator, repository Persister) *ShortenUseCase {
	return &ShortenUseCase{generator: generator, repository: repository}
}

// ShortenURL takes a long url and returns a short url.
func (useCase ShortenUseCase) ShortenURL(ctx context.Context, url LongURL) (*ShortURL, error) {
	if err := urlIsInvalid(url); err != nil {
		return nil, err
	}

	code, err := useCase.generator.GenerateCode(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to shorten url")
	}

	shortURL := &ShortURL{raw: url.String(), code: code}

	if err := useCase.repository.Persist(ctx, shortURL); err != nil {
		return nil, errors.Wrap(err, "failed to persist url")
	}

	return shortURL, nil
}
