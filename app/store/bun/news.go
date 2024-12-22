package bun

import (
	"context"
	"grimoire/app/log"
	"grimoire/app/model"
)

func (r *SaitRepo) CreateNews(ctx context.Context, entry *model.DBNews) (*model.DBNews, error) {
	_, err := r.db.
		NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		r.logger.Error("store.SaitRepo.GetNewsSlice",
			log.String("error", err.Error()),
			log.Object("entry", entry),
		)
		return nil, err
	}

	err = r.db.
		NewSelect().
		Model(entry).
		Relation("Profile").
		Where("id = ?", entry.ID).
		Scan(ctx)
	if err != nil {
		r.logger.Error("store.SaitRepo.AddComment (SELECT)",
			log.String("error", err.Error()),
			log.Object("entry", entry),
		)
		return nil, err
	}
	return entry, nil
}

func (r *SaitRepo) GetNewsByID(ctx context.Context, id int) (*model.DBNews, error) {
	entry := &model.DBNews{}
	err := r.db.
		NewSelect().
		Model(entry).
		Relation("Profile").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		r.logger.Error("store.SaitRepo.GetNewsByID",
			log.String("error", err.Error()),
			log.Object("entry", entry),
		)
		return nil, err
	}

	var comments model.DBCommentSlice
	err = r.db.
		NewSelect().
		Model(&comments).
		Where("news_id = ?", entry.ID).
		Relation("Profile").
		OrderExpr("like_count DESC, created_at DESC").
		Scan(ctx)
	if err != nil {
		r.logger.Error("store.SaitRepo.GetComments",
			log.String("error", err.Error()),
			log.Int("news_id", entry.ID),
		)
	}

	entry.Comments = comments

	return entry, nil
}

func (r *SaitRepo) GetNewsSlice(ctx context.Context) (*model.DBNewsSlice, error) {
	entry := &model.DBNewsSlice{}
	err := r.db.
		NewSelect().
		Model(entry).
		Relation("Profile").
		OrderExpr("created_at DESC").
		Scan(ctx)
	if err != nil {
		r.logger.Error("store.SaitRepo.GetNewsSlice",
			log.String("error", err.Error()),
			log.Object("entry", entry),
		)
		return nil, err
	}
	return entry, nil
}

func (r *SaitRepo) UpdateNews(ctx context.Context, entry *model.DBNews) error {
	err := r.db.
		NewInsert().
		Model(entry).
		Scan(ctx)
	if err != nil {
		r.logger.Error("store.SaitRepo.UpdateNews",
			log.String("error", err.Error()),
			log.Object("entry", entry),
		)
		return err
	}
	return nil
}

func (r *SaitRepo) DeleteNews(ctx context.Context, id int) error {
	_, err := r.db.
		NewDelete().
		Model((*model.DBNews)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		r.logger.Error("store.SaitRepo.DeleteNews",
			log.String("error", err.Error()),
			log.Int("news_id", id),
		)
		return err
	}
	return nil
}
