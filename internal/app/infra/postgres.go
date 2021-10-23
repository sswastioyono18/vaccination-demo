package infra

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/sswastioyono18/vaccination-demo/config"

	_ "github.com/lib/pq"
)

// NewPostgreDatabase return gorp dbmap object with postgre options param
func NewPostgreDatabase(config *config.AppConfig) (*gorp.DbMap, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Name, config.DB.Pass))
	if err != nil {
		return nil, err
	}

	gorp := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	return gorp, nil
}
