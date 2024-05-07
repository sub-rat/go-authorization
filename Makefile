.PHONY: start build


APP_NAME		= go-rbac-boilerplate
BUILD_ROOT		= build

all: start

fmt:
	@find . -name "*.go" -type f -not -path "./vendor/*"|xargs gofmt -s -w

build:
	@go build -ldflags "-w -s" -o $(BUILD_ROOT)/$(APP_NAME)

start:
	@go run ./main.go api-server --config=./config/config.yaml --casbin_model=./config/casbin_model.conf

migrate:
	@go run ./main.go migrate --config=./config/config.yaml

setup:
	@go run ./main.go setup --config=./config/config.yaml --menu=./config/menu.yaml

swagger:
	@swag init --parseDependency --parseInternal -g internal/routes/swagger_route.go

clean:
	@rm -rf $(BUILD_ROOT)

up:
	docker-compose -f deployments/compose.yaml up -d

down:
	docker-compose -f deployments/compose.yaml down

restart:
	docker-compose -f deployments/compose.yaml restart

gen:
	@go run ./main.go gen -c ./config/config.yaml -n $(name)
