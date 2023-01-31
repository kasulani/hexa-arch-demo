package app

import (
	"context"

	"github.com/teris-io/shortid"
)

type (
	generator struct{}
)

// newGenerator returns an instance of generator.
func newGenerator() *generator {
	return &generator{}
}

// GenerateCode creates a unique short code.
func (g *generator) GenerateCode(_ context.Context) (string, error) {
	id, err := shortid.Generate()
	if err != nil {
		return "", err
	}
	return id, nil
}
