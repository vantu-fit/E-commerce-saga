package db_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
)

// TestUpdateProductInventoryTx test function for UpdateProductInventoryTx
func TestUpdateProductInventoryTx(t *testing.T) {
	uuidAccount := uuid.New()

	idempotencyKey := uuid.New()

	purchasedProducts := []db.PurchasedProduct{}

	// create category
	category, err := testStore.CreateCategory(context.Background(), db.CreateCategoryParams{
		Name:        "category",
		Description: "description",
	})
	require.NoError(t, err)
	require.NotEmpty(t, category)

	// create product
	for i := range 5 {
		product, err := testStore.CreateProduct(context.Background(), db.CreateProductParams{
			IDCategory:  category.ID,
			IDAccount:   uuidAccount,
			Name:        "product",
			Price:       1000,
			Description: "description",
			BrandName:   "brand",
			Inventory:   10,
		})
		require.NoError(t, err)
		require.NotEmpty(t, product)

		purchasedProducts = append(purchasedProducts, db.PurchasedProduct{
			ProductID: product.ID,
			Quantity:  product.Inventory - int64(i+1),
		})
	}

	err = testStore.UpdateProductInventoryTx(context.Background(), idempotencyKey, &purchasedProducts)
	require.NoError(t, err)
}

func TestRollbackProductInventoryTx(t *testing.T) {
	uuidAccount := uuid.New()

	idempotencyKey := uuid.New()

	purchasedProducts := []db.PurchasedProduct{}

	// create category
	category, err := testStore.CreateCategory(context.Background(), db.CreateCategoryParams{
		Name:        "category",
		Description: "description",
	})
	require.NoError(t, err)
	require.NotEmpty(t, category)

	// create product
	for i := range 5 {
		product, err := testStore.CreateProduct(context.Background(), db.CreateProductParams{
			IDCategory:  category.ID,
			IDAccount:   uuidAccount,
			Name:        "product",
			Price:       1000,
			Description: "description",
			BrandName:   "brand",
			Inventory:   10,
		})
		require.NoError(t, err)
		require.NotEmpty(t, product)

		purchasedProducts = append(purchasedProducts, db.PurchasedProduct{
			ProductID: product.ID,
			Quantity:  int64(i+1),
		})
	}

	err = testStore.UpdateProductInventoryTx(context.Background(), idempotencyKey, &purchasedProducts)
	require.NoError(t, err)
	
	err = testStore.RollbackProductInventoryTx(context.Background(), idempotencyKey, &purchasedProducts)
	require.NoError(t, err)
}
