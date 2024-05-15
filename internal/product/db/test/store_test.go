package db_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/product/service/command"
	"github.com/vantu-fit/saga-pattern/internal/product/service/entity"
)

// TestUpdateProductInventoryTx test function for UpdateProductInventoryTx
func TestUpdateProductInventoryTx(t *testing.T) {

	// create category
	category, err := testStore.CreateCategory(context.Background(), db.CreateCategoryParams{
		Name:        "category",
		Description: "description",
	})
	require.NoError(t, err)
	require.NotEmpty(t, category)

	purchasedProducts := []db.PurchasedProduct{}
	// create product
	for i:= 0 ; i < 5; i++ {
		product, err := testStore.CreateProduct(context.Background(), db.CreateProductParams{
			IDCategory:  category.ID,
			IDAccount:   uuid.New(),
			Name:        "product",
			Price:       1000,
			Description: "description",
			BrandName:   "brand",
			Inventory:   1000,
		})

		require.NoError(t, err)
		require.NotEmpty(t, product)

		purchasedProducts = append(purchasedProducts, db.PurchasedProduct{
			ProductID: product.ID,
			Quantity:  int64(10),
		})
	}
	n :=  64
	wg := sync.WaitGroup{}
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			id := uuid.New()
			err = testStore.UpdateProductInventoryTx(context.Background(), id, &purchasedProducts)
			errs <- err
		}(i , &wg)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	for _, purchasedProduct := range purchasedProducts {
		product, err := testStore.GetProductForUpdate(context.Background(), purchasedProduct.ProductID)
		require.NoError(t, err)
		require.NotEmpty(t, product)
		require.Equal(t, int64(1000-10 * n), product.Inventory)
	}

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
			Quantity:  int64(i + 1),
		})
	}

	err = testStore.UpdateProductInventoryTx(context.Background(), idempotencyKey, &purchasedProducts)
	require.NoError(t, err)

	err = testStore.RollbackProductInventoryTx(context.Background(), idempotencyKey, &purchasedProducts)
	require.NoError(t, err)
}

func TestUpdateProductInventoryTxCommand(t *testing.T) {
	updateProductInventoryHanlder := command.NewUpdateProductInventoryHandler(testStore)

	purchasedProducts := make([]entity.ProductItem, 5)

	// create category
	category, err := testStore.CreateCategory(context.Background(), db.CreateCategoryParams{
		Name:        "category",
		Description: "description",
	})
	require.NoError(t, err)
	require.NotEmpty(t, category)

	// create product
	for i  := range 5 {
		product, err := testStore.CreateProduct(context.Background(), db.CreateProductParams{
			IDCategory:  category.ID,
			IDAccount:   uuid.New(),
			Name:        "product",
			Price:       1000,
			Description: "description",
			BrandName:   "brand",
			Inventory:   10000,
		})

		require.NoError(t, err)
		require.NotEmpty(t, product)
		fmt.Println(product.ID)
		purchasedProducts[i] =  entity.ProductItem{
			ID:       product.ID,
			Quantity: int64(10),
		}
	}


	wg := sync.WaitGroup{}
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = updateProductInventoryHanlder.Handle(context.Background(), command.UpdateProductInventory{
				PurchaseID:   uuid.New(),
				ProductItems: &purchasedProducts,
			})
			require.NoError(t, err)
		}()
	}
	wg.Wait()

}
