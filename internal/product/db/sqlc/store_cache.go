package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/cache"
	"github.com/vantu-fit/saga-pattern/pkg/utils"
)

const (
	cuckooFilter = "product_cuckoo_filter"
	dummnyItem   = "dummy_item"
	mutexKey     = "mutex:"

	checkProductKey        = "check_product:"
	getProductDetailKey    = "get_product_detail:"
	getProductInventoryKey = "get_product_inventory:"
	getProductKey          = "get_product:"
)

type CacheStore struct {
	*SQLStore
	store *SQLStore
	lc    cache.LocalCache
	rc    cache.RedisCache
}

func NewCacheStore(store *SQLStore, lc cache.LocalCache, rc cache.RedisCache) Store {
	return &CacheStore{
		SQLStore: store,
		store:    store,
		lc:       lc,
		rc:       rc,
	}
}

func (storeCache *CacheStore) UpdateProductInventoryTx(ctx context.Context, idempotencyKey uuid.UUID, purchasedProducts *[]PurchasedProduct) error {
	err := storeCache.store.UpdateProductInventoryTx(ctx, idempotencyKey, purchasedProducts)
	if err != nil {
		return err
	}

	payloads := make([]cache.RedisIncrbyXPayload, len(*purchasedProducts))
	for i, purchasedProduct := range *purchasedProducts {
		payloads[i] = cache.RedisIncrbyXPayload{
			Key:   utils.StrJoin(getProductInventoryKey, purchasedProduct.ProductID.String()),
			Value: int64(-purchasedProduct.Quantity),
		}
	}
	if len(payloads) > 0 {
		err = storeCache.rc.ExecIncrbyXPipeline(ctx, &payloads)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (storeCache *CacheStore) RollbackProductInventoryTx(ctx context.Context, idempotencyKey uuid.UUID, purchasedProducts *[]PurchasedProduct) error {
	err := storeCache.store.RollbackProductInventoryTx(ctx, idempotencyKey, purchasedProducts)
	if err != nil {
		return err
	}

	payloads := make([]cache.RedisIncrbyXPayload, len(*purchasedProducts))
	for i, purchasedProduct := range *purchasedProducts {
		payloads[i] = cache.RedisIncrbyXPayload{
			Key:   utils.StrJoin(getProductInventoryKey, purchasedProduct.ProductID.String()),
			Value: int64(purchasedProduct.Quantity),
		}
	}
	if len(payloads) > 0 {
		err = storeCache.rc.ExecIncrbyXPipeline(ctx, &payloads)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
