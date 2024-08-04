/*
リクエストが切断されるとContextがCancel状態になり、DBとの通信とかも終了されてしまう。
非同期処理をリクエスト終了後も続けるために、Contextをクローンしたい
*/
package context

import (
	"context"
	"time"
)

// キャンセルされずに裏で処理を続けたいとき用
/*
	ctx, _ := AsyncContext(ctx)

	go func() {
		dbClient.CreateCard(ctx, card)
	}
*/
func AsyncContext(ctx context.Context) (context.Context, context.CancelFunc) {
	asyncContext := context.WithoutCancel(ctx)
	asyncContext, cancel := context.WithTimeout(asyncContext, 15*time.Second)
	return asyncContext, cancel
}
