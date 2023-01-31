package url

import "context"

type (
	mockCodeGenerator func(ctx context.Context) (string, error)
	mockPersister     func(ctx context.Context, url *ShortURL) (err error)
)

// GenerateCode is a mock.
func (m mockCodeGenerator) GenerateCode(ctx context.Context) (string, error) {
	return m(ctx)
}

// Persist is a mock.
func (m mockPersister) Persist(ctx context.Context, url *ShortURL) (err error) {
	return m(ctx, url)
}
