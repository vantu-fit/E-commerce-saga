ACCOUNT_DB_URL=postgresql://postgres:postgres@localhost:5432/account_db?sslmode=disable
PRODUCT_DB_URL=postgresql://postgres:postgres@localhost:5432/product_db?sslmode=disable
ORDER_DB_URL=postgresql://postgres:postgres@localhost:5432/order_db?sslmode=disable


protoc:
	rm -f pb/*.go 
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	proto/*.proto
redis:
	docker compose -f ./docker-compose-redis.yml up
kafka:
	docker compose -f ./docker-compose-kafka.yml up

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


test:
	pg_ctl start
	make redis
	make account
	make product


	