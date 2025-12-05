#!/bin/bash

# Run all Python asyncio examples
# 运行所有 Python asyncio 示例

echo "=========================================="
echo "Python Asyncio 示例"
echo "Python Asyncio Examples"
echo "=========================================="

examples=(
    "01-asyncio-basics"
    "02-event-loop"
    "03-streams-and-protocols"
    "04-practical-examples"
    "05-custom-implementation"
    "06-advanced-examples"
)

for example in "${examples[@]}"; do
    echo ""
    echo "=========================================="
    echo "运行示例: $example"
    echo "Running: $example"
    echo "=========================================="
    
    cd "python-asyncio/$example" || exit
    
    # Run all .py files in the directory
    for file in *.py; do
        if [ -f "$file" ]; then
            echo "运行 $file..."
            python3 "$file"
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
