


run:
	@echo '== Generating code =='
	go generate ./...
	@echo '== Running compose =='
	docker-compose up -d
	@echo '== Running app =='
	APP_DEV=true DB_HOST=localhost KAFKA_HOST=localhost go run cmd/main.go