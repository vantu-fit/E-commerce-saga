ACCOUNT_DB_URL=postgresql://postgres:postgres@localhost:5432/account_db?sslmode=disable
PRODUCT_DB_URL=postgresql://postgres:postgres@localhost:5432/product_db?sslmode=disable
ORDER_DB_URL=postgresql://postgres:postgres@localhost:5432/order_db?sslmode=disable
PAYMENT_DB_URL=postgresql://postgres:postgres@localhost:5432/payment_db?sslmode=disable
MEDIA_DB_URL=postgresql://postgres:postgres@localhost:5432/media_db?sslmode=disable


protoc:
	rm -f pb/*.go 
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	proto/*.proto
redis:
	docker compose -f ./docker-compose-redis.yml up
redis-down:
	docker compose -f ./docker-compose-redis.yml down
kafka:
	docker compose -f ./docker-compose-kafka.yml up
minio:
	docker compose -f ./docker-compose-minio.yml up
minio-down:
	docker compose -f ./docker-compose-minio.yml down
envoy:
	docker compose -f ./docker-compose-envoy.yml up

account:
	go run ./cmd/account/main.go
mock-account:
	mockgen -package mockdb -destination internal/account/db/mock/store.go github.com/vantu-fit/saga-pattern/internal/account/db/sqlc Store
migrate-account-up:
	migrate -path internal/account/db/migration -database "$(ACCOUNT_DB_URL)" -verbose up
migrate-account-down:
	migrate -path internal/account/db/migration -database "$(ACCOUNT_DB_URL)" -verbose down
sqlc-account:
	sqlc generate -f sqlc/account.yml



product:
	go run ./cmd/product/main.go
migrate-product-up:
	migrate -path internal/product/db/migration -database "$(PRODUCT_DB_URL)" -verbose up
migrate-product-down:
	migrate -path internal/product/db/migration -database "$(PRODUCT_DB_URL)" -verbose down
sqlc-product:
	sqlc generate -f sqlc/product.yml


order:
	go run ./cmd/order/main.go
migrate-order-up:
	migrate -path internal/order/db/migration -database "$(ORDER_DB_URL)" -verbose up
migrate-order-down:
	migrate -path internal/order/db/migration -database "$(ORDER_DB_URL)" -verbose down
sqlc-order:
	sqlc generate -f sqlc/order.yml


payment:
	go run ./cmd/payment/main.go
migrate-payment-up:
	migrate -path internal/payment/db/migration -database "$(PAYMENT_DB_URL)" -verbose up
migrate-payment-down:
	migrate -path internal/payment/db/migration -database "$(PAYMENT_DB_URL)" -verbose down
sqlc-payment:
	sqlc generate -f sqlc/payment.yml

purchase:
	go run ./cmd/purchase/main.go


orchestrator:
	go run ./cmd/orchestrator/main.go

media:
	go run ./cmd/media/main.go
migrate-media-up:
	migrate -path internal/media/db/migration -database "$(MEDIA_DB_URL)" -verbose up
migrate-media-down:
	migrate -path internal/media/db/migration -database "$(MEDIA_DB_URL)" -verbose down
sqlc-media:
	sqlc generate -f sqlc/media.yml




test:
	pg_ctl start
	make redis
	make account
	make product


	