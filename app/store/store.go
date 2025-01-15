package store

import (
	"context"
	"warhoop/app/config"
	"warhoop/app/log"
	"warhoop/app/store/bun"

	_ "github.com/go-sql-driver/mysql"
)

// Store ...
type Store struct {
	AuthRepo *bun.AuthRepo
	SaitRepo *bun.SaitRepo
	CharRepo *bun.CharRepo
}

func NewBun(ctx context.Context, logger *log.Logger) (*Store, error) {
	cfg := config.Get()
	if cfg == nil {
		panic("configuration is nil")
	}

	auth, err := connect(ctx, cfg.DB.Auth, cfg.DB.Verbose, logger, "auth")
	if err != nil {
		logger.Error("Failed to initialize auth",
			log.String("error", err.Error()),
		)
		return nil, err
	}

	char, err := connect(ctx, cfg.DB.Characters, cfg.DB.Verbose, logger, "characters")
	if err != nil {
		logger.Error("Failed to initialize characters",
			log.String("error", err.Error()),
		)
		return nil, err
	}

	sait, err := connect(ctx, cfg.DB.Sait, cfg.DB.Verbose, logger, "sait")
	if err != nil {
		logger.Error("Failed to initialize sait",
			log.String("error", err.Error()),
		)
		return nil, err
	}

	saitRepo := bun.NewSaitRepo(sait, logger)
	if saitRepo == nil {
		logger.Error("Failed to initialize SaitRepo",
			log.String("error", err.Error()),
		)
		return nil, err
	}

	authRepo := bun.NewAuthRepo(auth, logger, saitRepo)
	if authRepo == nil {
		if err != nil {
			logger.Error("Failed to initialize AuthRepo",
				log.String("error", err.Error()),
			)
			return nil, err
		}
	}

	charRepo := bun.NewCharRepo(char, logger)
	if charRepo == nil {
		if err != nil {
			logger.Error("Failed to initialize AuthRepo",
				log.String("error", err.Error()),
			)
			return nil, err
		}
	}
	return &Store{
		AuthRepo: authRepo,
		CharRepo: charRepo,
		SaitRepo: saitRepo,
	}, nil
}

func connect(ctx context.Context, dsn string, verbose bool, logger *log.Logger, name string) (*bun.DB, error) {
	db, err := bun.Dial(dsn, verbose)
	if err != nil {
		logger.Error("Dial failed",
			log.String("db_name", name),
			log.String("err", err.Error()),
		)
		return nil, err
	}

	tracerName := "my-application-tracer"

	db.AddQueryHook(bun.NewOtelQueryHook(tracerName, name))

	return db, nil
}
