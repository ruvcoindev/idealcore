.PHONY: test test-hypercube test-web test-short

# Тестировать всё (может падать)
test:
	go test ./pkg/... -v

# Тестировать только hypercube
test-hypercube:
	go test ./pkg/hypercube/... -v

# Тестировать hypercube с покрытием
test-hypercube-coverage:
	go test ./pkg/hypercube/... -coverprofile=hypercube.out
	go tool cover -html=hypercube.out -o hypercube.html
	@echo "📊 Coverage:"
	go tool cover -func=hypercube.out | grep total

# Тестировать web (работает)
test-web:
	go test ./pkg/web -v

# Быстрое тестирование работающих пакетов
test-short:
	go test ./pkg/hypercube/... -short
	go test ./pkg/web -short

# Починить зависимости
tidy:
	go mod tidy
