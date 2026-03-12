#!/bin/bash

echo "🔧 Исправление импортов в pkg/hypercube..."

find pkg/hypercube -name "*.go" -type f -exec sed -i \
  's|github.com/ruvcoindev/idealcore/hypercube/|github.com/ruvcoindev/idealcore/pkg/hypercube/|g' {} +

find pkg/hypercube -name "*.go" -type f -exec sed -i \
  's|github.com/idealcore/hypercube/|github.com/ruvcoindev/idealcore/pkg/hypercube/|g' {} +

echo "✅ Исправления завершены"
