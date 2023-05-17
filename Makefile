PROTO_FILES=$(shell find proto/v1 -name '*.proto')
MIGRATIONS_DIR=internal/app/repository/postgresql/migrations
DOCKER_COMPOSE_FILE=build/docker-compose.yaml

generate-grpc-service:
	protoc \
	--go_out=. \
	--go-grpc_out=. \
	--grpc-gateway_out=. \
	--grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out . \
	$(PROTO_FILES)

	# protoc generates swagger files near by proto definitions, let's move them to separate directory "openapiv2" to better structure
	build/move_openapiv2_specs.sh "proto/v1" "proto/v1/openapiv2/"

services-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up

services-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

migrate-up:
	goose -dir=$(MIGRATIONS_DIR) postgres $(DB_URL) up

migrate-down:
	goose -dir=$(MIGRATIONS_DIR) postgres $(DB_URL) down

export-env:
	set -o allexport && source .env && set +o allexport
