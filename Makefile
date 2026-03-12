.PHONY: test test-coverage test-race bench lint clean

# Запуск всех тестов
test:
	go test ./pkg/... -v

# Запуск тестов с покрытием
test-coverage:
	go test ./pkg/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Запуск тестов с проверкой гонок
test-race:
	go test ./pkg/... -race -v

# Запуск бенчмарков
bench:
	go test ./pkg/... -bench=. -benchmem

# Запуск линтера
lint:
	golangci-lint run ./pkg/...

# Очистка
clean:
	rm -f coverage.out coverage.html
	rm -f *.log

# Запуск всех проверок
all: test test-race lint

# Запуск тестов для конкретного пакета
test-core:
	go test ./pkg/hypercube/core -v

test-data:
	go test ./pkg/hypercube/data -v

test-model:
	go test ./pkg/hypercube/model -v

test-analysis:
	go test ./pkg/hypercube/analysis -v

test-protocol:
	go test ./pkg/hypercube/protocol -v

# Запуск тестов с детальным выводом
test-verbose:
	go test ./pkg/... -v -cover

# Запуск тестов с коротким выводом
test-short:
	go test ./pkg/... -short

# Запуск тестов с таймаутом
test-timeout:
	go test ./pkg/... -timeout 30s

help:
	@echo "Available commands:"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-race      - Run tests with race detection"
	@echo "  make bench          - Run benchmarks"
	@echo "  make lint           - Run linter"
	@echo "  make clean          - Clean generated files"
	@echo "  make all            - Run all checks"
	@echo "  make test-core      - Test core package only"
	@echo "  make test-data      - Test data package only"
	@echo "  make test-model     - Test model package only"
	@echo "  make test-analysis  - Test analysis package only"
	@echo "  make test-protocol  - Test protocol package only"
