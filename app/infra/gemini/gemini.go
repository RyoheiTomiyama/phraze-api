package gemini

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain/infra/gemini"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type client struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

type ClientOption struct {
	APIKey string
}

func New(opts ClientOption) (gemini.IClient, error) {
	ak := option.WithAPIKey(opts.APIKey)

	cl, err := genai.NewClient(context.Background(), ak)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	model := cl.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	return &client{
		client: cl,
		model:  model,
	}, nil
}
