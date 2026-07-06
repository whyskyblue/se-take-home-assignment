#!/bin/bash

# Run Script
# This script should execute your CLI application and output results to result.txt

echo "Running CLI application..."

cd "$(dirname "$0")/.."

chmod +x ./order-controller

./order-controller <<EOF
+
n
v
n
v
+
s
q
EOF

if [ $? -eq 0 ]; then
    echo "CLI application execution completed"
else
    echo "CLI application execution failed"
    exit 1
fi
