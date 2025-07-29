# LeetCode MCP Server

A Model Context Protocol (MCP) server written in Go that integrates with VS Code to fetch daily LeetCode challenges.

## Overview

This MCP server provides a tool to retrieve the current LeetCode daily challenge, including the problem title, difficulty, description, and direct link. It's designed to work seamlessly with VS Code through the MCP protocol.

## Features

- 🎯 **Daily Challenge Retrieval**: Get today's LeetCode daily challenge problem
- 📝 **Rich Content**: Includes problem description, difficulty, and formatted content
- 🔗 **Direct Links**: Provides direct links to the problem on LeetCode
- 🧹 **Clean Formatting**: Converts HTML content to readable markdown-like text
- ⚡ **Fast Response**: Efficient GraphQL API integration with proper error handling

## Installation

### Prerequisites

- Go 1.19 or higher
- VS Code with MCP support

### Building the Server

```bash
# Clone or navigate to the project directory
cd leetcode-mcp-go

# Build the server
go build -o leetcode-mcp-server main.go
```

## Configuration

### VS Code Integration

To use this MCP server with VS Code, you need to configure it in your MCP client settings. Add the following configuration:

```json
{
  "mcpServers": {
    "leetcode": {
      "command": "/path/to/leetcode-mcp-server",
      "args": []
    }
  }
}
```

## Usage

Once configured and running, the server provides the following tool:

### github

uses oauth2 to authenticate with GitHub

```
prompt 1: on mcp-workshop repo, create issue "feature: add dark mode"
prompt 2: on mcp-workshop repo, create issue "feature: add light mode", generate description text as bullet list on why light mode is needed
prompt 3: on mentioned repo, create a python flask api for products in a PR, assign agent to work on it
```

### playwright

```text
prompt 1: navigate to https://tfl.gov.uk/ I'm going from paddington to heathrow
prompt 2: generate playwright tests based on above interaction

```

### apple docs

```text
Card swiping gestures for a card UI #apple-developer-docs
```

### `get_leetcode_daily_challenge`


```text
show me daily challenge, use tool
```

Fetches today's LeetCode daily challenge problem.

**Parameters:**
- `include_content` (boolean, optional): Whether to include the full problem description (default: true)

**Returns:**
- Problem title
- Difficulty level
- Problem ID
- Direct link to the problem
- Problem description (if requested)
- Date of the challenge

**Example response:**
```
**Two Sum**
**Difficulty:** Easy
**Problem ID:** 1
**Link:** https://leetcode.com/problems/two-sum
**Date:** 2025-07-22

**Problem Description:**
Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target...
```

## Development

### Running the Server

```bash
# Run directly with Go
go run main.go

# Or build and run the binary
go build -o leetcode-mcp-server main.go
./leetcode-mcp-server
```

The server communicates via stdin/stdout using JSON-RPC 2.0 protocol as specified by MCP.

### Testing

You can test the server by sending JSON-RPC requests via stdin:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./leetcode-mcp-server
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./leetcode-mcp-server
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_leetcode_daily_challenge","arguments":{}}}' | ./leetcode-mcp-server
```

## Architecture

### MCP Protocol Implementation

The server implements the core MCP methods:

- **`initialize`**: Establishes the protocol version and server capabilities
- **`tools/list`**: Returns available tools
- **`tools/call`**: Executes the requested tool

### LeetCode API Integration

Uses LeetCode's GraphQL API to fetch daily challenge data:
- Endpoint: `https://leetcode.com/graphql`
- Query: `activeDailyCodingChallengeQuestion`
- Content formatting: Converts HTML to readable text

## Error Handling

The server includes comprehensive error handling:
- Network request failures
- JSON parsing errors
- Invalid parameters
- API response errors
- Proper JSON-RPC error codes

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source. Feel free to use and modify as needed.

## Related Links

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [LeetCode](https://leetcode.com/)
- [VS Code MCP Integration](https://code.visualstudio.com/)
