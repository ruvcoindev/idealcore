#!/bin/bash
echo "🔍 Тестирование hypercube..."

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Функция для тестирования пакета
test_package() {
    echo -e "\n📦 Тестирование $1..."
    if go test "./pkg/hypercube/$1" -v; then
        echo -e "${GREEN}✅ $1 OK${NC}"
        return 0
    else
        echo -e "${RED}❌ $1 FAILED${NC}"
        return 1
    fi
}

# Тестируем по порядку
test_package "core"
test_package "data"
test_package "model"
test_package "analysis"
test_package "protocol"

# Интеграционный тест
echo -e "\n🔄 Интеграционное тестирование..."
if go test "./pkg/hypercube" -v; then
    echo -e "${GREEN}✅ Интеграция OK${NC}"
else
    echo -e "${RED}❌ Интеграция FAILED${NC}"
fi

echo -e "\n📊 Покрытие кода:"
go test ./pkg/hypercube/... -coverprofile=hypercube.out
go tool cover -func=hypercube.out | grep total
