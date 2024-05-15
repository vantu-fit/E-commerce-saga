package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/service/entity"
)

var (
	HandlerHeader = "handler"

	// PurchaseTopic is the subscribed topic for new purchase
	PurchaseTopic       = "purchase"
	PurchaseGroupID     = "purchase-group"
	PurchaseResultTopic = "purchase-result"

	// UpdateProductInventoryTopic is the topic to which we publish update product inventory
	UpdateProductInventoryTopic   = "update-product-inventory"
	UpdateProductInventoryGroupID = "update-product-inventory-group"
	UpdateProductInventoryHandler = "update-product-inventory-handler"

	RollbackProductInventoryTopic   = "rollback-product-inventory"
	RollbackProductInventoryGroupID = "rollback-product-inventory-group"
	RollbackProductInventoryHandler = "rollback-product-inventory-handler"

	// CreateOrderTopic is the topic to which we publish create order
	CreateOrderTopic   = "create-order"
	CreateOrderGroupID = "create-order-group"
	CreateOrderHandler = "create-order-handler"

	RollbackOrderTopic   = "rollback-order"
	RollbackOrderGroupID = "rollback-order-group"
	RollbackOrderHandler = "rollback-order-handler"

	// CreatePaymentTopic is the topic to which we publish create order
	CreatePaymentTopic   = "create-payment"
	CreatePaymentGroupID = "create-payment-group"
	CreatePaymentHandler = "create-payment-handler"

	RollbackPaymentTopic   = "rollback-payment"
	RollbackPaymentGroupID = "rollback-payment-group"
	RollbackPaymentHandler = "rollback-payment-handler"

	// ReplyTopic is saga step reply topic
	ReplyTopic   = "reply"
	ReplyGroupID = "reply-group"

	// Mail
	SendRegisterEmailTopic   = "send-register-email"
	SendRegisterEmailGroupID = "send-register-email-group"
	SendRegisterEmailHandler = "send-register-email-handler"
)

var (
	StepUpdateProductInventory = "UPDATE_PRODUCT_INVENTORY"
	StepCreateOrder            = "CREATE_ORDER"
	StepCreatePayment          = "CREATE_PAYMENT"

	StatusExecute        = "EXUCUTE"
	StatusSuccess        = "SUCCESS"
	StatusFailed         = "FAILED"
	StatusRollback       = "ROLLBACK"
	StatusRollbackFailed = "ROLLBACK_FAILED"
)

// PurchaseResult event
type PurchaseResult struct {
	PurchaseID uuid.UUID
	Step       string
	Status     string
	Timestamp  time.Time
}

type CreatePurchaseResponse struct {
	Purchase *entity.Purchase
	Succsess bool
	Error    string
}
