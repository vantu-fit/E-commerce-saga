package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	UpdateProductInventoryTx(ctx context.Context, idempotencyKey uuid.UUID, purchasedProducts *[]PurchasedProduct) error
	RollbackProductInventoryTx(ctx context.Context, idempotencyKey uuid.UUID, purchasedProducts *[]PurchasedProduct) error
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
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
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
	// time.Sleep(5 * time.Second)
	return tx.Commit(ctx)
}

type PurchasedProduct struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int32     `json:"quantity"`
}

func (store *SQLStore) UpdateProductInventoryTx(ctx context.Context, idempotencyKey uuid.UUID, purchasedProducts *[]PurchasedProduct) error {
	return store.execTx(ctx, func(q *Queries) error {
		// check idempotency key
		idempotencies, err := store.GetIdempotencyKey(ctx, idempotencyKey)
		if err != nil {
			return err
		}

		if len(idempotencies) > 0 {
			return fmt.Errorf("idempotency key %v already exists", idempotencyKey)
		}

		// update purchase products inventory
		for _, purchasedProduct := range *purchasedProducts {
			// get product for update : lock
			product, err := store.GetProductForUpdate(ctx, purchasedProduct.ProductID)
			if err != nil {
				return err
			}

			// chekc inventory
			if product.Inventory < purchasedProduct.Quantity {
				return fmt.Errorf("product %v has insufficient inventory", purchasedProduct.ProductID)
			}

			argUpadateProduct := UpadateProductParams{
				ID: product.ID,
				Inventory: pgtype.Int4{
					Int32: product.Inventory - purchasedProduct.Quantity,
					Valid: true,
				},
			}

			// update product inventory
			_, err = store.UpadateProduct(ctx, argUpadateProduct)
			if err != nil {
				return err
			}
		}

		// insert idempotency key
		for _, purchasedProduct := range *purchasedProducts {
			_, err := store.CreateIdempotency(ctx, CreateIdempotencyParams{
				ID:         idempotencyKey,
				ProductID:  purchasedProduct.ProductID,
				Quantity:   purchasedProduct.Quantity,
				Rollbacked: false,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (store *SQLStore) RollbackProductInventoryTx(ctx context.Context, idempotencyKey uuid.UUID, purchasedProducts *[]PurchasedProduct) error {
	return store.execTx(ctx, func(q *Queries) error {
		// get idempotencies key of this purchase
		idempotencies, err := store.GetIdempotencyKey(ctx, idempotencyKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				return fmt.Errorf("idempotency key %v not found", idempotencyKey)
			}
			return err
		}

		// update purchase products inventory
		for _, idempotency := range idempotencies {
			// get product for update : lock
			product, err := store.GetProductForUpdate(ctx, idempotency.ProductID)
			if err != nil {
				return err
			}

			// check rollbacked
			if idempotency.Rollbacked {
				continue
			}

			argUpadateProduct := UpadateProductParams{
				ID: product.ID,
				Inventory: pgtype.Int4{
					Int32: product.Inventory + idempotency.Quantity,
					Valid: true,
				},
			}

			// update product inventory
			_, err = store.UpadateProduct(ctx, argUpadateProduct)
			if err != nil {
				return err
			}
		}

		// update idempotency
		_, err = store.UpdateIdempotency(ctx, UpdateIdempotencyParams{
			ID: idempotencyKey,
			Rollbacked: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})

		return nil
	})
}
