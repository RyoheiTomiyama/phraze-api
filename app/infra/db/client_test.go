package db

import (
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
)

func testDataSourceOption() DataSourceOption {
	config, err := env.New()
	if err != nil {
		return DataSourceOption{}
	}

	return DataSourceOption{
		Host:     config.DB.HOST,
		Port:     config.DB.PORT,
		DBName:   db_test.TestDBName,
		User:     config.DB.USER,
		Password: config.DB.PASSWORD,
	}
}
