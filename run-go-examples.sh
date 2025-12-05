#!/bin/bash

# Run all Go async examples
# 运行所有 Go 异步示例

echo "=========================================="
echo "Go 异步编程示例"
echo "Go Async Programming Examples"
echo "=========================================="

examples=(
    "01-basic-goroutines"
    "02-channels"
    "03-select-statement"
    "04-context"
    "05-sync-package"
    "06-advanced-patterns"
    "07-practical-examples"
)

for example in "${examples[@]}"; do
    echo ""
    echo "=========================================="
    echo "运行示例: $example"
    echo "Running: $example"
    echo "=========================================="
    
    cd "go-async/$example" || exit
    
    # Run all .go files in the directory
    for file in *.go; do
        if [ -f "$file" ]; then
            echo "运行 $file..."
            go run "$file"
            echo ""
        fi
    done
    
    cd ../..
done

echo ""
echo "=========================================="
echo "所有示例运行完成！"
echo "All examples completed!"
echo "=========================================="
