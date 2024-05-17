package db

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vantu-fit/saga-pattern/pb"
)

type Store interface {
	Querier
	DeleteCommentTx(ctx context.Context, arg DeleteCommentParamsTx) error
	CreateCommentTx(ctx context.Context, arg CreateCommentParamsTx) error
}
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{
		// IsoLevel: pgx.ReadUncommitted,
	})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type CreateCommentParamsTx struct {
	*pb.CreateCommentRequest
	UserID uuid.UUID
}

// CreateComment Transaction
func (store *SQLStore) CreateCommentTx(ctx context.Context, arg CreateCommentParamsTx) error {
	return store.execTx(ctx, func(q *Queries) error {
		// check parent comment
		if arg.ParentId == nil {
			_, err := store.Queries.CreateComment(ctx, CreateCommentParams{
				ProductID:  uuid.MustParse(arg.ProductId),
				UserID:     arg.UserID,
				Content:    arg.Content,
				LeftIndex:  1,
				RightIndex: 2,
			})
			if err != nil {
				return err
			}
			return nil
		}

		// get parent comment for update
		parentComment, err := store.Queries.GetCommentForUpdate(ctx, uuid.MustParse(*arg.ParentId))
		if err != nil {
			return err
		}
		// get max right index
		rightIndex := parentComment.RightIndex

		// check is root comment
		if !parentComment.ParentID.Valid {
			// update left index
			_, err = store.Queries.UpdateLeftIndexComment(ctx, UpdateLeftIndexCommentParams{
				ParentID: pgtype.UUID{
					Bytes: parentComment.ID,
					Valid: true,
				},
				LeftIndex:   rightIndex,
				LeftIndex_2: 2,
			})
			if err != nil {
				return err
			}

			// update right index
			_, err = store.Queries.UpdateRightIndexComment(ctx, UpdateRightIndexCommentParams{
				ParentID: pgtype.UUID{
					Bytes: parentComment.ID,
					Valid: true,
				},
				RightIndex:   rightIndex,
				RightIndex_2: 2,
			})
			if err != nil {
				return err
			}

			// create comment
			_, err = store.Queries.CreateComment(ctx, CreateCommentParams{
				ProductID:  uuid.MustParse(arg.ProductId),
				UserID:     arg.UserID,
				Content:    arg.Content,
				LeftIndex:  rightIndex,
				RightIndex: rightIndex + 1,
				ParentID: pgtype.UUID{
					Bytes: parentComment.ID,
					Valid: true,
				},
			})
			if err != nil {
				return err
			}

			return nil

		}

		// update left index
		_, err = store.Queries.UpdateLeftIndexComment(ctx, UpdateLeftIndexCommentParams{
			ParentID:    parentComment.ParentID,
			LeftIndex:   rightIndex,
			LeftIndex_2: 2,
		})
		if err != nil {
			return err
		}

		// update right index
		_, err = store.Queries.UpdateRightIndexComment(ctx, UpdateRightIndexCommentParams{
			ParentID:     parentComment.ParentID,
			RightIndex:   rightIndex,
			RightIndex_2: 2,
		})
		if err != nil {
			return err
		}

		// create comment
		_, err = store.Queries.CreateComment(ctx, CreateCommentParams{
			ProductID:  uuid.MustParse(arg.ProductId),
			UserID:     arg.UserID,
			Content:    arg.Content,
			LeftIndex:  rightIndex,
			RightIndex: rightIndex + 1,
			ParentID:   parentComment.ParentID,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

type DeleteCommentParamsTx struct {
	*pb.DeleteCommentRequest
}

// DeleteComment Transaction
func (store *SQLStore) DeleteCommentTx(ctx context.Context, arg DeleteCommentParamsTx) error {
	return store.execTx(ctx, func(q *Queries) error {
		// get comment for update
		comment, err := store.Queries.GetCommentForUpdate(ctx, uuid.MustParse(arg.Id))
		if err != nil {
			log.Println(err)
			return err
		}


		// check is root comment
		if !comment.ParentID.Valid {
			// delete comment
			_, err = store.Queries.DeleteComment(ctx, DeleteCommentParams{
				ParentID: pgtype.UUID{
					Bytes: comment.ID,
					Valid: true,
				},
				LeftIndex:  comment.LeftIndex,
				RightIndex: comment.RightIndex,
			})
			if err != nil {
				return err
			}

			return nil
		}

		width := comment.RightIndex - comment.LeftIndex + 1

		// delete comment
		_, err = store.Queries.DeleteComment(ctx, DeleteCommentParams{
			ParentID:   comment.ParentID,
			LeftIndex:  comment.LeftIndex,
			RightIndex: comment.RightIndex,
		})
		if err != nil {
			return err
		}

		// update left index
		_, err = store.Queries.UpdateLeftIndexComment(ctx, UpdateLeftIndexCommentParams{
			ParentID:    comment.ParentID,
			LeftIndex:   comment.RightIndex,
			LeftIndex_2: -width,
		})
		if err != nil {
			return err
		}

		// update right index
		_, err = store.Queries.UpdateRightIndexComment(ctx, UpdateRightIndexCommentParams{
			ParentID:     comment.ParentID,
			RightIndex:   comment.RightIndex,
			RightIndex_2: -width,
		})
		if err != nil {
			return err
		}

		return nil
	})
}
