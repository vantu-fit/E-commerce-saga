package db

// import (
// 	"context"
// 	"fmt"

// 	"github.com/google/uuid"
// 	"github.com/rs/zerolog/log"
// 	"github.com/vantu-fit/saga-pattern/pkg/cache"
// 	"github.com/vantu-fit/saga-pattern/pkg/utils"
// )

// const (
// 	cuckooFilter = "product_cuckoo_filter"
// 	dummnyItem   = "dummy_item"
// 	mutexKey     = "mutex:"

// 	getProductInventoryKey = "get_product_inventory:"
// 	getProductKey          = "get_product:"
// )

// type StoreCache struct {
// 	Store
// 	lc cache.LocalCache
// 	rc cache.RedisCache
// }

// func NewStoreCache(store Store, lc cache.LocalCache, rc cache.RedisCache) Store {
// 	exist, err := rc.CFExist(context.Background(), cuckooFilter, dummnyItem)
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to check cuckoo filter existence, err: %s", err)
// 		return store
// 	}

// 	if !exist {
// 		err = rc.CFReserve(context.Background(), cuckooFilter, 1000, 4, 1000)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to reserve cuckoo filter, err: %s", err)
// 			return store
// 		}

// 		err = rc.CFAdd(context.Background(), cuckooFilter, dummnyItem)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to add dummy item to cuckoo filter, err: %s", err)
// 			return store
// 		}

// 		log.Info().Msg("Product: reserved cuckoo filter")
// 	}

// 	return &StoreCache{
// 		Store: store,
// 		lc:    lc,
// 		rc:    rc,
// 	}
// }

// func (sc *StoreCache) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
// 	product, err := sc.Store.CreateProduct(ctx, arg)
// 	if err != nil {
// 		return product, err
// 	}

// 	err = sc.rc.CFAdd(ctx, cuckooFilter, product.ID.String())
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to add product to cuckoo filter, err: %s", err)
// 	}

// 	return product, nil
// }

// func (sc *StoreCache) DeleteProduct(ctx context.Context, id uuid.UUID) (Product, error) {
// 	exist, err := sc.rc.CFExist(ctx, cuckooFilter, id.String())
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to check product existence, err: %s", err)
// 		return Product{}, err
// 	}

// 	if exist {
// 		err = sc.rc.CFDel(ctx, cuckooFilter, id.String())
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to delete product from cuckoo filter, err: %s", err)
// 		}
// 	}

// 	ok, err := sc.lc.Get(getProductKey+id.String(), &Product{})
// 	if ok && err == nil {
// 		err = sc.lc.Delete(getProductKey + id.String())
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to delete product from local cache, err: %s", err)
// 		}
// 	}

// 	ok, err = sc.lc.Get(getProductKey+id.String(), &Product{})
// 	if ok && err == nil {
// 		err = sc.lc.Delete(getProductKey + id.String())
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to delete product detail from local cache, err: %s", err)
// 		}
// 	}

// 	return sc.Store.DeleteProduct(ctx, id)
// }

// func (sc *StoreCache) GetProductInventory(ctx context.Context, id uuid.UUID) (GetProductInventoryRow, error) {
// 	key := utils.StrJoin(getProductInventoryKey, id.String())
// 	var inventoryProduct int64
// 	var inventory = GetProductInventoryRow{
// 		ID: 	  id,
// 	}

// 	ok, err := sc.lc.Get(key, &inventoryProduct)
// 	if ok && err == nil {
// 		inventory.Inventory = inventoryProduct
// 		return inventory, nil
// 	}

// 	exist, err := sc.rc.CFExist(ctx, cuckooFilter, id.String())
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to check product existence, err: %s", err)
// 		return GetProductInventoryRow{}, err
// 	}

// 	if !exist {
// 		inventory, err := sc.Store.GetProductInventory(ctx, id)
// 		if err != nil {
// 			return inventory, err
// 		}

// 		err = sc.rc.CFAdd(ctx, cuckooFilter, id.String())
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to add product to cuckoo filter, err: %s", err)
// 		}

// 		err = sc.lc.Set(key, &inventory.Inventory)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to set product inventory to local cache, err: %s", err)
// 		}

// 		err = sc.rc.Set(ctx, key, &inventory.Inventory)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to set product inventory to redis cache, err: %s", err)
// 		}

// 		return inventory, nil
// 	}

// 	ok, err = sc.rc.Get(ctx, key, &inventoryProduct)
// 	if ok && err == nil {
// 		err = sc.lc.Set(key, &inventoryProduct)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to set product inventory to local cache, err: %s", err)
// 		}
// 	}

// 	// lock to prevent lost update
// 	mu := sc.rc.GetMutex(mutexKey + key)
// 	err = mu.Lock()
// 	if err != nil {
// 		return inventory, err
// 	}
// 	defer mu.Unlock()

// 	// Get again to prevent new update
// 	ok, err = sc.rc.Get(ctx, key, &inventoryProduct)
// 	if ok && err == nil {
// 		err = sc.lc.Set(key, &inventoryProduct)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to set product inventory to local cache, err: %s", err)
// 		}
// 		inventory.Inventory = inventoryProduct
// 		return inventory, nil
// 	}

// 	inventory, err = sc.Store.GetProductInventory(ctx, id)
// 	if err != nil {
// 		return inventory, err
// 	}

// 	err = sc.rc.Set(ctx, key, &inventory.Inventory)
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to set product inventory to redis cache, err: %s", err)
// 	}

// 	err = sc.lc.Set(key, &inventory.Inventory)
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to set product inventory to local cache, err: %s", err)
// 	}

// 	return inventory, nil

// }

// func (sc *StoreCache) GetProductByID(ctx context.Context, id uuid.UUID) (Product, error) { 
// 	var product Product
// 	key := getProductKey + id.String()

// 	ok, err := sc.lc.Get(key, &product)
// 	if ok && err == nil {
// 		return product, nil
// 	}

// 	exist, err := sc.rc.CFExist(ctx, cuckooFilter, id.String())
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to check product existence, err: %s", err)
// 		return Product{}, err
// 	}

// 	if !exist {
// 		product , err := sc.Store.GetProductByID(ctx, id)
// 		if err != nil {
// 			return product, err
// 		}

// 		err = sc.rc.CFAdd(ctx, cuckooFilter, id.String())
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to add product to cuckoo filter, err: %s", err)
// 		}

// 		err = sc.lc.Set(key, &product)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to set product to local cache, err: %s", err)
// 		}

// 		err = sc.rc.Set(ctx, key, &product)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to set product to redis cache, err: %s", err)
// 		}

// 		return product, nil
// 	}

// 	ok, err = sc.rc.Get(ctx, key, &product)
// 	if ok && err == nil {
// 		err = sc.lc.Set(key, &product)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to set product to local cache, err: %s", err)
// 		}
// 	}

// 	// lock to prevent lost update
// 	mu := sc.rc.GetMutex(mutexKey + key)
// 	err = mu.Lock()
// 	if err != nil {
// 		return product, err
// 	}	
// 	defer mu.Unlock()

// 	// Get again to prevent new update
// 	ok, err = sc.rc.Get(ctx, key, &product)
// 	if ok && err == nil {
// 		err = sc.lc.Set(key, &product)
// 		if err != nil {
// 			log.Error().Msgf("Product: failed to set product to local cache, err: %s", err)
// 		}
// 		return product, nil
// 	}
	
// 	product, err = sc.Store.GetProductByID(ctx, id)
// 	if err != nil {
// 		return product, err
// 	}

// 	err = sc.rc.Set(ctx, key, &product)
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to set product to redis cache, err: %s", err)
// 	}

// 	err = sc.lc.Set(key, &product)
// 	if err != nil {
// 		log.Error().Msgf("Product: failed to set product to local cache, err: %s", err)
// 	}

// 	return product, nil
// }


// func (storeCache *StoreCache) UpdateProductInventoryTx(ctx context.Context, idempotencyKey uuid.UUID, purchasedProducts *[]PurchasedProduct) error {
// 	err := storeCache.Store.UpdateProductInventoryTx(ctx, idempotencyKey, purchasedProducts)
// 	if err != nil {
// 		log.Error().Msgf("UpdateProductInventoryTx: failed to update product inventory, err: %s", err)
// 		return err
// 	}

// 	payloads := make([]cache.RedisIncrbyXPayload, len(*purchasedProducts))
// 	for i, purchasedProduct := range *purchasedProducts {
// 		payloads[i] = cache.RedisIncrbyXPayload{
// 			Key:   utils.StrJoin(getProductInventoryKey, purchasedProduct.ProductID.String()),
// 			Value: int64(-purchasedProduct.Quantity),
// 		}
// 	}
// 	if len(payloads) > 0 {
// 		err = storeCache.rc.ExecIncrbyXPipeline(ctx, &payloads)
// 		if err != nil {
// 			log.Error().Msgf("UpdateProductInventoryTx: failed to update product inventory, err: %s", err)
// 			return err
// 		}
// 		return nil
// 	}

// 	return nil
// }

// func (storeCache *StoreCache) RollbackProductInventoryTx(ctx context.Context, idempotencyKey uuid.UUID, purchasedProducts *[]PurchasedProduct) error {
// 	err := storeCache.Store.RollbackProductInventoryTx(ctx, idempotencyKey, purchasedProducts)
// 	if err != nil {
// 		fmt.Println("RollbackProductInventoryTx", err)
// 		return err
// 	}

// 	payloads := make([]cache.RedisIncrbyXPayload, len(*purchasedProducts))
// 	for i, purchasedProduct := range *purchasedProducts {
// 		payloads[i] = cache.RedisIncrbyXPayload{
// 			Key:   utils.StrJoin(getProductInventoryKey, purchasedProduct.ProductID.String()),
// 			Value: int64(purchasedProduct.Quantity),
// 		}
// 	}
// 	if len(payloads) > 0 {
// 		err = storeCache.rc.ExecIncrbyXPipeline(ctx, &payloads)
// 		if err != nil {
// 			fmt.Println("RollbackProductInventoryTx", err)
// 			return err
// 		}
// 		return nil
// 	}

// 	return nil
// }
