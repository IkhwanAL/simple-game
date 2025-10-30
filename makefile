APP_NAME = tinyworlds
PKG = ./...

ifeq ($(OS),Windows_NT)
	BIN_NAME = $(APP_NAME).exe
else
	BIN_NAME = $(APP_NAME)
endif

run:
	@echo "🏃 Running $(APP_NAME) with race detection..."
	templ generate
	npx @tailwindcss/cli -i ./assets/input.css -o ./static/tailwind.css
	go run -race cmd/server/main.go

generate:
	@echo "🧩 Generating templ + tailwind..."
	templ generate
	npx @tailwindcss/cli -i ./assets/input.css -o ./static/tailwind.css

tests:
	@echo "🧪 Running tests with race detection..."
	go test -race -v $(PKG)

lint:
	@echo "🔍 Linting with staticcheck..."
	staticcheck $(PKG)

fmt:
	go fmt $(PKG)

build:
	@echo "Im Building A Go Binary"
	templ generate
	npx @tailwindcss/cli -i ./assets/input.css -o ./static/tailwind.css 
	go build -race -o bin/$(BIN_NAME) cmd/server/main.go
