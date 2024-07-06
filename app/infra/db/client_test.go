package db

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func NewTestClient(t *testing.T, db *sqlx.DB) IClient {
	return &client{db}
}
