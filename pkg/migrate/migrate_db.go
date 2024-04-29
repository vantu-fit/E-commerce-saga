package migrate_db

import (
	"github.com/rs/zerolog/log"
	"github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/github"


	
)

func RunDBMigration(migrationUrl string, dbSource string) {
	migraton, err := migrate.New(
		migrationUrl,
		dbSource,
	)

	if err != nil {
		log.Fatal().Msgf("cannot create migrate instance: %s", err)
		return
	}

	if err := migraton.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("cannot migrate db: %s", err)
	}

	log.Print("db migrate successfully")
}
