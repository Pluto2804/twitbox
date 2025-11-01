#############################################
# PROJECT CONFIG
#############################################

main_package_path = ./cmd/web
binary_name = example

#############################################
# LOAD ENV FILE
#############################################

ifneq (,$(wildcard .env))
include .env
export $(shell sed 's/=.*//' .env)
endif

#############################################
# HELPERS
#############################################

## help: Show help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$$(git status --porcelain)"

#############################################
# QUALITY CONTROL
#############################################

## audit: Run quality checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$$(gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: Run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: Run tests with coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## upgradeable: List upgrades for modules
.PHONY: upgradeable
upgradeable:
	@go run github.com/oligot/go-mod-upgrade@latest

#############################################
# DEVELOPMENT
#############################################

## tidy: Format and tidy
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## build: Build binary
.PHONY: build
build:
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

## run: Run app (dev)
.PHONY: run
run:
	go run ${main_package_path}

## run/live: Auto reload dev mode
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" \
		--build.bin "/tmp/bin/${binary_name}" \
		--build.delay "100" \
		--build.include_ext "go,tpl,tmpl,html,css,scss,js,ts,sql,jpg,jpeg,gif,png,svg,ico"

#############################################
# OPERATIONS
#############################################

## push: Push to Git
.PHONY: push
push: confirm audit no-dirty
	git push

## production/deploy: Build for Linux + compress
#############################################
# PRODUCTION DEPLOYMENT
#############################################

# Your server IP
production_host_ip = 143.244.138.123
production_user = twitbox
binary_name = example
main_package_path = ./cmd/web

## production/connect: SSH into the server
.PHONY: production/connect
production/connect:
	ssh $(production_user)@$(production_host_ip)


## production/deploy: Build, upload & migrate
.PHONY: production/deploy
production/deploy: production/build
	@echo ">> Copying binary and migrations to server..."
	rsync -P /tmp/bin/$(binary_name) $(production_user)@$(production_host_ip):~
	rsync -rp --delete ./migrations $(production_user)@$(production_host_ip):~

	@echo ">> Running migrations on server..."
	ssh -t $(production_user)@$(production_host_ip) \
		"bash -l -c 'migrate -path ~/migrations -database \"$$DB_DSN\" up'"




	@echo "✅ Deployment complete"


## production/build: Build Linux binary
.PHONY: production/build
production/build:
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/$(binary_name) $(main_package_path)



#############################################
# DATABASE & MIGRATIONS
#############################################

## mysql: Open MySQL CLI
#############################################
# DATABASE & MIGRATIONS
#############################################

## mysql: Open MySQL CLI
.PHONY: mysql
mysql:
	mysql -u$(DB_USER) -p$(DB_PASS) -h$(DB_HOST) -P$(DB_PORT) $(DB_NAME)

## migrate/up: Apply migrations
.PHONY: migrate/up
migrate/up:
	@echo "Running migrations..."
	migrate -path ./migrations -database "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true" up
	@echo "Migrations completed ✅"

## migrate/down: Rollback last migration
.PHONY: migrate/down
migrate/down:
	@echo "Rolling back migrations..."
	migrate -path ./migrations -database "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true" down
	@echo "Rollback complete ✅"
