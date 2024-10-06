package gemini

import "context"

type IClient interface {
	GenAnswer(ctx context.Context, q string) (string, error)
}
