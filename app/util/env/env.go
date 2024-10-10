// 環境変数をContextに詰める
// FromCtxから利用可能
package env

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	APP    app
	DB     db
	Gemini gemini
	Sentry sentry
}

type db struct {
	HOST     string `env:"POSTGRES_HOST" envDefault:"0.0.0.0"`
	USER     string `env:"POSTGRES_USER" envDefault:"postgres"`
	PASSWORD string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	DB_NAME  string `env:"POSTGRES_DB" envDefault:"phraze"`
	PORT     string `env:"POSTGRES_PORT" envDefault:"5432"`
}

type gemini struct {
	API_KEY string `env:"GOOGLE_GEMINI_API_SECRET_KEY"`
}

type sentry struct {
	DSN string `env:"SENTRY_DSN"`
}

type app struct {
	CORS []string `env:"CORS" envSeparator:"," envDefault:"http://localhost:3000"`
	ENV  string   `env:"ENV"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, errutil.Wrap((err))
	}

	return &cfg, nil
}

type envCtxKey struct{}

func (c *Config) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, envCtxKey{}, c)
}

func FromCtx(ctx context.Context) *Config {
	c, ok := ctx.Value(envCtxKey{}).(*Config)
	if !ok {
		return &Config{}
	}

	return c
}
