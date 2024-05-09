package service

import (
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/product/service/command"
	"github.com/vantu-fit/saga-pattern/internal/product/service/query"
)

type Command struct {
	CreateProduct            command.CreateProductHandler
	CreateCategory           command.CreateCategoryHandler
	UpdateProductInventory   command.UpdateProductInventoryHandler
	RollbackProductInventory command.RollbackProductInventoryHandler
	UpadateProductDetail     command.UpdateProductDetailHandler
}

type Query struct {
	ListProduct  query.ListProductHandler
	GetProduct   query.GetProductHandler
	CheckProduct query.CheckProductHandler
}

type Service struct {
	Command Command
	Query   Query
}

func NewService(store db.Store) Service {
	return Service{
		Command: Command{
			CreateProduct:            command.NewCreateProductHandler(store),
			CreateCategory:           command.NewCreateCategoryHandler(store),
			UpdateProductInventory:   command.NewUpdateProductInventoryHandler(store),
			RollbackProductInventory: command.NewRollbackProductInventoryHandler(store),
			UpadateProductDetail:     command.NewUpdateProductDetailHandler(store),
		},
		Query: Query{
			ListProduct:  query.NewListProductHandler(store),
			GetProduct:   query.NewGetProductHandler(store),
			CheckProduct: query.NewCheckProductHandler(store),
		},
	}
}
