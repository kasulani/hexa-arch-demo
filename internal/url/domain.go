package url

import (
	"context"
	"regexp"
	"strings"
)

type (
	// ShortURL is an entity.
	ShortURL struct {
		raw  string
		code string
	}

	// LongURL is a value object.
	LongURL string

	// CodeGenerator is a port for a unique short code generator.
	CodeGenerator interface {
		GenerateCode(ctx context.Context) (string, error)
	}

	// Shortener is a port for the url shorten use case.
	Shortener interface {
		ShortenURL(ctx context.Context, url LongURL) (*ShortURL, error)
	}

	// Persister is a port for saving the short url to a repository.
	Persister interface {
		Persist(ctx context.Context, url *ShortURL) (err error)
	}

	// ErrInvalidURL is a custom error for an invalid url.
	ErrInvalidURL struct{}
)

// String implements fmt.Stringer interface.
func (s LongURL) String() string {
	return string(s)
}

// Code returns the short url code.
func (s *ShortURL) Code() string {
	return s.code
}

// Raw returns the raw url.
func (s *ShortURL) Raw() string {
	return s.raw
}

// Error implements the builtin error interface.
func (e ErrInvalidURL) Error() string {
	return "invalid url"
}

// urlIsInvalid is an example of an invariant. In DDD, an invariant is business rule.
func urlIsInvalid(url LongURL) error {
	exp, err := regexp.Compile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)`)
	if err != nil {
		return err
	}

	value := url.String()
	if i := strings.Index(value, "#"); i > -1 {
		value = value[:i]
	}

	if !exp.MatchString(value) {
		return ErrInvalidURL{}
	}

	return nil
}
