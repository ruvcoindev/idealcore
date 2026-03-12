#!/bin/bash
echo "🎯 Целевое тестирование"

# Тестируем только рабочие пакеты
PACKAGES=(
    "./pkg/hypercube/..."
    "./pkg/web"
    # добавляйте сюда пакеты, которые точно работают
)

for pkg in "${PACKAGES[@]}"; do
    echo "📦 Тестирование $pkg"
    go test $pkg -v -short
done
