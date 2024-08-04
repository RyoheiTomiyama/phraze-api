package context

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type traceIDKey struct{}

func TestAsyncContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), traceIDKey{}, "req123")

	t.Run("ContextのValueが継承されること", func(t *testing.T) {
		clonedCtx, _ := AsyncContext(ctx)
		assert.Equal(t, "req123", clonedCtx.Value(traceIDKey{}))
	})

	for _, tc := range []struct {
		Name      string
		Arrange   func(ctx context.Context) context.Context
		AsyncFunc func(something *int)
		Assert    func(result int)
	}{
		{
			Name: "ContextがCancelしたら非同期処理が停止する例の確認",
			Arrange: func(ctx context.Context) context.Context {
				return ctx
			},
			AsyncFunc: func(something *int) {
				*something = 1
			},
			Assert: func(result int) {
				assert.Equal(t, 0, result)
			},
		},
		{
			Name: "親ContextがCancelしたら非同期処理が停止する例の確認",
			Arrange: func(ctx context.Context) context.Context {
				type AnyKey string
				ctx = context.WithValue(ctx, AnyKey("any"), "any")
				return ctx
			},
			AsyncFunc: func(something *int) {
				*something = 1
			},
			Assert: func(result int) {
				assert.Equal(t, 0, result)
			},
		},
		{
			Name: "CloneContextの場合、親ContextがCancelしても非同期処理が停止しないこと",
			Arrange: func(ctx context.Context) context.Context {
				ctx, _ = AsyncContext(ctx)
				return ctx
			},
			AsyncFunc: func(something *int) {
				*something = 1
			},
			Assert: func(result int) {
				assert.Equal(t, 1, result)
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)

			ctx = tc.Arrange(ctx)

			something := 0
			var wg sync.WaitGroup

			// Contextがcancelしたら停止する処理
			asyncTestLogic(t, ctx, &wg, func() {
				tc.AsyncFunc(&something)
				t.Log("something async logic")
			})

			// contextをcancelしてリクエストを終了した状態を作る
			cancel()
			t.Log("original context canceled")

			wg.Wait()

			tc.Assert(something)
		})
	}
}

// contextがcancelされていた場合fnを発火させないメソッド
func asyncTestLogic(t *testing.T, ctx context.Context, wg *sync.WaitGroup, fn func()) {
	t.Helper()

	done := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			t.Log("context done")
		case <-done:
			t.Log("done!")
		}
	}()

	go func() {
		// 非同期処理
		time.Sleep(500 * time.Millisecond)

		// Contextがcancelされていたら処理停止する
		select {
		case <-ctx.Done():
			return
		default:
			fn()
			close(done)
		}
	}()
}
