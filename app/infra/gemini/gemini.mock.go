package gemini

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockedClient struct {
	mock.Mock
}

func NewMock() *MockedClient {
	return &MockedClient{}
}

func (m *MockedClient) GenAnswer(ctx context.Context, q string) (string, error) {
	args := m.Called(ctx, q)
	return args.String(0), args.Error(1)
}
