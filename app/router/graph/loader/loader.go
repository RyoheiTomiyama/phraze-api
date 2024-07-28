package loader

import (
	"context"
)

// dataloadgenのloaderのインターフェイス
// https://github.com/vikstrous/dataloadgen/blob/v0.0.6/dataloadgen.go
type LoaderInterface[KeyT comparable, ValueT any] interface {
	Clear(key KeyT)
	Load(ctx context.Context, key KeyT) (ValueT, error)
	LoadAll(ctx context.Context, keys []KeyT) ([]ValueT, error)
	LoadAllThunk(ctx context.Context, keys []KeyT) func() ([]ValueT, error)
	LoadThunk(ctx context.Context, key KeyT) func() (ValueT, error)
	Prime(key KeyT, value ValueT) bool
}
