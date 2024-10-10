package gemini

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/util/env"
	"github.com/stretchr/testify/assert"
)

func TestGenAnswer(t *testing.T) {
	// Geminiのトークン消費してしまうのでテスト省略
	// 動作確認したいときは手元でSkipをコメントアウトして実行する
	t.Skip()

	ctx := context.Background()
	config, err := env.New()
	if err != nil {
		panic(err)
	}
	ctx = config.WithCtx(ctx)

	c, err := New(ClientOption{APIKey: config.Gemini.API_KEY})
	if err != nil {
		t.Fatal(err)
	}

	ans, err := c.GenAnswer(ctx, "answer")
	assert.NoError(t, err)
	t.Log(ans)
}
