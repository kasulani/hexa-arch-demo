package url

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestShortenUseCase(t *testing.T) {
	testURL := LongURL("https://www.testing.com")

	tests := map[string]struct {
		url              LongURL
		generator        CodeGenerator
		persist          Persister
		expectedShortURL *ShortURL
		expectedErr      error
	}{
		"when the unique code generator fails, return an error": {
			expectedErr: errors.New("failed to shorten url: a-generator-error"),
			url:         testURL,
			generator: mockCodeGenerator(
				func(ctx context.Context) (string, error) {
					return "", errors.New("a-generator-error")
				},
			),
		},
		"when the persister fails, return an error": {
			expectedErr: errors.New("failed to persist url: a-database-error"),
			url:         testURL,
			generator: mockCodeGenerator(
				func(ctx context.Context) (string, error) {
					return "xyz", nil
				},
			),
			persist: mockPersister(
				func(ctx context.Context, url *ShortURL) (err error) {
					return errors.New("a-database-error")
				},
			),
		},
		"when the long url is invalid, return an error": {
			expectedErr: ErrInvalidURL{},
			url:         "https://www.invalid-url",
		},
		"happy path": {
			url: testURL,
			generator: mockCodeGenerator(
				func(ctx context.Context) (string, error) {
					return "abCDE", nil
				},
			),
			persist: mockPersister(
				func(ctx context.Context, url *ShortURL) (err error) {
					return nil
				},
			),
			expectedShortURL: &ShortURL{
				raw:  "https://www.test.com",
				code: "abCDE",
			},
		},
	}

	for name, test := range tests {
		testCase := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			usecase := NewShortenURLUseCase(testCase.generator, testCase.persist)
			shorturl, err := usecase.ShortenURL(context.Background(), testCase.url)

			switch testCase.expectedErr != nil {
			case true:
				assert.EqualError(t, err, testCase.expectedErr.Error())
				assert.Empty(t, shorturl)
			case false:
				assert.NoError(t, err)
				assert.NotEmpty(t, shorturl)
			}
		})
	}
}
