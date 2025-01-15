package bun

import (
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
	"warhoop/app/log"
)

type DB struct {
	*bun.DB
}

func Dial(dbDSN string, verbose bool) (*DB, error) {
	if dbDSN == "" {
		err := errors.New("missing DSN: database connection string is empty")
		log.Get().Error("Dial failed",
			log.String("error", err.Error()),
		)
		return nil, err
	}

	sqlDB, err := sql.Open("mysql", dbDSN)
	if err != nil {
		log.Get().Error("Dial: failed to open database connection",
			log.String("error", err.Error()),
		)
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Get().Error("Dial: failed to ping database",
			log.String("error", err.Error()),
		)
		return nil, err
	}

	bunDB := bun.NewDB(sqlDB, mysqldialect.New(), bun.WithDiscardUnknownColumns())

	bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(verbose)))

	_, err = bunDB.Exec("SELECT 1")
	if err != nil {
		log.Get().Error("Dial: failed to execute test query",
			log.String("error", err.Error()),
		)
		return nil, err
	}

	log.Get().Debug("Dial: database connection established successfully",
		log.String("DSN", dbDSN),
	)

	return &DB{bunDB}, nil
}
