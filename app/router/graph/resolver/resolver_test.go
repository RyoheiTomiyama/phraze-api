package resolver_test

import (
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/card"
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/deck"
	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/resolver"
	card_service "github.com/RyoheiTomiyama/phraze-api/service/card"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type resolverSuite struct {
	suite.Suite
	resolver *resolver.Resolver
	dbx      *sqlx.DB
	dbClient db.IClient
}

func TestResolverSuite(t *testing.T) {
	suite.Run(t, &resolverSuite{})
}

// テスト間で干渉してしまうので、テストごとにtxdbをリセットする
func (s *resolverSuite) SetupTest() {
	dbx := db_test.GetDB(s.T())
	dbClient := db.NewTestClient(s.T(), dbx)
	cardService := card_service.NewService()

	cardUsecase := card.New(dbClient, cardService)
	deckUsecase := deck.New(dbClient)

	resolver := resolver.New(cardUsecase, deckUsecase)

	s.resolver = resolver
	s.dbx = dbx
	s.dbClient = dbClient
}

func (s *resolverSuite) TearDownTest() {
	s.dbx.Close()
}
