up:
	test -f .env || cp .env.example .env
	docker-compose up --build

down:
	docker-compose down -v --remove-orphans

migrate:
	docker run -v ./db/migrations:/migrations \
		--network test-go-project_app-network \
		migrate/migrate \
		-path=/migrations \
		-database postgres://postgres:pwd@db:5432/master?sslmode=disable \
		up

lint:
	golangci-lint run -v

test:
	go test -v ./...

fmt:
	go fmt ./... && gofumpt -w .

codegen:
	oapi-codegen -generate types -o internal/api/dto.go -package generated api/openapi.yaml
