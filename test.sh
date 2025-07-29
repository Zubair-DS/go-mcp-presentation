#!/bin/bash

# Test script for LeetCode MCP Server

echo "Building the MCP server..."
go build -o leetcode-mcp-server main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"
echo ""

echo "Testing MCP server initialization..."
response=$(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./leetcode-mcp-server | head -1)
echo "Response: $response"
echo ""

echo "Testing tools list..."
response=$(echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./leetcode-mcp-server | head -1)
echo "Response: $response"
echo ""

echo "Testing daily challenge retrieval..."
echo "This may take a few seconds to fetch from LeetCode API..."
response=$((echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}'; echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_leetcode_daily_challenge","arguments":{}}}') | ./leetcode-mcp-server | tail -1)
echo "Response: $response"
echo ""

echo "✅ All tests completed!"
echo ""
echo "To use this MCP server with VS Code:"
echo "1. Make sure you have MCP support in VS Code"
echo "2. Configure the server path in your MCP client settings"
echo "3. Use the 'get_leetcode_daily_challenge' tool to fetch daily challenges"
