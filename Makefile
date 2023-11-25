# Constants

PROJECT_NAME=payment-service
USER=fidesy-pay


PHONY: generate
generate:
	mkdir -p pkg/${PROJECT_NAME}
	protoc --go_out=pkg/${PROJECT_NAME} --go_opt=paths=import \
			--go-grpc_out=pkg/${PROJECT_NAME} --go-grpc_opt=paths=import \
			--grpc-gateway_out=pkg/${PROJECT_NAME} \
            --grpc-gateway_opt grpc_api_configuration=./api/${PROJECT_NAME}/${PROJECT_NAME}.yaml \
            --grpc-gateway_opt allow_delete_body=true \
			api/${PROJECT_NAME}/${PROJECT_NAME}.proto
	mv pkg/${PROJECT_NAME}/github.com/${USER}/${PROJECT_NAME}/* pkg/${PROJECT_NAME}
	rm -r pkg/${PROJECT_NAME}/github.com

PHONY: go-build
go-build:
	GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/${PROJECT_NAME}
	mkdir -p bin
	mv main bin

PHONY: build
build:
	make go-build
	docker build --tag ${PROJECT_NAME} .

PHONY: run
run:
	docker run --name ${PROJECT_NAME} -dp 7777:7777 -e GRPC_PORT=7777 ${PROJECT_NAME}

PHONY: migrate-up
migrate-up:
	goose -dir ./migrations postgres "postgresql://user:pass@host:port/db?sslmode=disable" up

PHONY: migrate-down
migrate-down:
	goose -dir ./migrations postgres "postgresql://user:pass@host:port/db?sslmode=disable" down

PHONY: generate-swagger
generate-swagger:
	protoc -I . --openapiv2_out ./ \
	  --experimental_allow_proto3_optional=true \
      --openapiv2_opt grpc_api_configuration=./api/$(PROJECT_NAME)/$(PROJECT_NAME).yaml \
      --openapiv2_opt proto3_optional_nullable=true \
      --openapiv2_opt allow_delete_body=true \
      ./api/$(PROJECT_NAME)/$(PROJECT_NAME).proto

	mv api/$(PROJECT_NAME)/$(PROJECT_NAME).swagger.json ./swaggerui/swagger_temp.json
	jq '. + {"host": "$(APP_NAME).fidesy.xyz:$(PROXY_PORT)", "schemes": ["http"]}' ./swaggerui/swagger_temp.json > ./swaggerui/swagger.json
	rm ./swaggerui/swagger_temp.json
