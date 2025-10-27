APP_NAME = tinyworlds
PKG = ./...

run:
	@echo "🏃 Running $(APP_NAME) with race detection..."
	go run -race cmd/main.go

generate:
	@echo "🧩 Generating templ + tailwind..."
	templ generate
	npx @tailwindcss/cli -i ./input.css -o ./static/tailwind.css --watch

test:
	@echo "🧪 Running tests with race detection..."
	go test -race -v $(PKG)

lint:
	@echo "🔍 Linting with staticcheck..."
	staticcheck $(PKG)

fmt:
	go fmt $(PKG)

build:
	go build -o bin/$(APP_NAME) cmd/main.go
