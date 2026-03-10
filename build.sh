#!/bin/bash
set -e

echo "📦 idealcore — сборка"

# Зависимости
go mod tidy

# Pure-Go версия
echo "🔨 Сборка Pure-Go..."
CGO_ENABLED=0 go build -o idealcore ./cmd/server/main.go

# CGO версия (если есть hnswlib)
if [ -f /usr/local/include/hnswlib/hnswalg.h ]; then
    echo "🔨 Сборка CGO + HNSW..."
    CGO_ENABLED=1 go build -tags cgo -o idealcore_cgo ./cmd/server/main.go
fi

echo "✅ Готово!"
ls -lh idealcore*
