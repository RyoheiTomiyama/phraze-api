package resolver

import (
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/deck"
	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type resolverSuite struct {
	suite.Suite
	resolver *Resolver
	dbx      *sqlx.DB
}

func TestResolverSuite(t *testing.T) {
	dbx := db_test.GetDB(t)
	defer dbx.Close()

	dbClient := db.NewTestClient(t, dbx)

	deckUsecase := deck.New(dbClient)

	resolver := &Resolver{
		deckUsecase: deckUsecase,
	}

	suite.Run(t, &resolverSuite{resolver: resolver, dbx: dbx})
}
