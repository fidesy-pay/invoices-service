# Constants

PROJECT_NAME=invoices-service
USER=fidesy-pay
APP_NAME=${PROJECT_NAME}-stage

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

PHONY: clean
clean:
	 if docker inspect ${APP_NAME} > /dev/null 2>&1; then docker rm -f ${APP_NAME} && docker rmi -f ${APP_NAME}; else echo "Container not found."; fi

PHONY: go-build
go-build:
	GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/${PROJECT_NAME}
	mkdir -p bin
	mv main bin

PHONY: build
build:
	make go-build
	docker build --tag ${APP_NAME} .

PHONY: run
run:
	make clean
	make build
	docker run --name ${APP_NAME} --network=zoo -dp 7030:7030 -e GRPC_PORT=7030 -e PROXY_PORT=7031 -e SWAGGER_PORT=7032 -e METRICS_PORT=7033 -e APP_NAME=${APP_NAME} -e ENV=local ${APP_NAME}

PHONY: migrate-up
migrate-up:
	goose -dir ./migrations postgres "postgres://invoices-service:invoices-service@localhost:45001/invoices-service?sslmode=disable" up

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
	jq '. + {"host": "$(SERVER_HOST):$(PROXY_PORT)", "schemes": ["http"]}' ./swaggerui/swagger_temp.json > ./swaggerui/swagger.json
	rm ./swaggerui/swagger_temp.json
