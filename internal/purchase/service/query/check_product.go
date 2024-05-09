package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/purchase/service/entity"
)

type OrderItem struct {
	ID       uuid.UUID
	Quantity int64
}

type CheckProducts struct {
	OrderItems *[]OrderItem
}

type CheckProductHandler QueryHandler[CheckProducts , *[]entity.ProductStatus]

type checkProductHandkler struct { 

}


func NewCheckProductHandler() CheckProductHandler {
	return &checkProductHandkler{

	}
}

func (h *checkProductHandkler) Handle (ctx context.Context , cmd CheckProducts) (*[]entity.ProductStatus , error){
	return nil , nil
}