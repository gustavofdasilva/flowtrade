proto:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

dev:
	go run services/account-service/main.go \
	go run services/api-gateway/main.go