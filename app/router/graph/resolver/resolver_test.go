package resolver

import (
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/card"
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/deck"
	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	card_service "github.com/RyoheiTomiyama/phraze-api/service/card"
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
	cardService := card_service.NewService()

	cardUsecase := card.New(dbClient, cardService)
	deckUsecase := deck.New(dbClient)

	resolver := &Resolver{
		cardUsecase,
		deckUsecase,
	}

	suite.Run(t, &resolverSuite{resolver: resolver, dbx: dbx})
}
