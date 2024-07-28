package loader

import (
	"context"

	"github.com/vikstrous/dataloadgen"
)

type noCacheLoader[KeyT comparable, ValueT any] struct {
	loader *dataloadgen.Loader[KeyT, ValueT]
}

// キャッシュを再利用しないLoader
// dataloadgenのローダーは同じキーが飛んできても、キャッシュから利用してしまいsqlを叩かないので、
// 例えばユーザーの継承ロールが更新されたあとも、以前の継承ロールを返してしまう。
// Primeを利用して適宜更新を行えばキャッシュを利用できるが、Prime漏れが発生したら障害になるので、キャッシュを利用しないようにしている
func NewNoCacheLoader[KeyT comparable, ValueT any](
	fetch func(ctx context.Context, keys []KeyT) ([]ValueT, []error),
	options ...dataloadgen.Option,
) LoaderInterface[KeyT, ValueT] {

	loader := dataloadgen.NewLoader(fetch, options...)

	return &noCacheLoader[KeyT, ValueT]{loader}
}

func (nl *noCacheLoader[KeyT, ValueT]) Clear(key KeyT) {
	nl.loader.Clear(key)
}

func (nl *noCacheLoader[KeyT, ValueT]) clears(keys []KeyT) {
	for _, key := range keys {
		nl.Clear(key)
	}
}

func (nl *noCacheLoader[KeyT, ValueT]) Load(ctx context.Context, key KeyT) (ValueT, error) {
	value, err := nl.loader.Load(ctx, key)
	nl.Clear(key)

	return value, err
}

func (nl *noCacheLoader[KeyT, ValueT]) LoadAll(ctx context.Context, keys []KeyT) ([]ValueT, error) {
	values, err := nl.loader.LoadAll(ctx, keys)
	nl.clears(keys)

	return values, err
}

func (nl *noCacheLoader[KeyT, ValueT]) LoadAllThunk(ctx context.Context, keys []KeyT) func() ([]ValueT, error) {
	thunk := nl.loader.LoadAllThunk(ctx, keys)

	return func() ([]ValueT, error) {
		values, err := thunk()
		nl.clears(keys)

		return values, err
	}
}

func (nl *noCacheLoader[KeyT, ValueT]) LoadThunk(ctx context.Context, key KeyT) func() (ValueT, error) {
	thunk := nl.loader.LoadThunk(ctx, key)

	return func() (ValueT, error) {
		value, err := thunk()
		nl.Clear(key)

		return value, err
	}
}

func (nl *noCacheLoader[KeyT, ValueT]) Prime(key KeyT, value ValueT) bool {
	return nl.loader.Prime(key, value)
}
