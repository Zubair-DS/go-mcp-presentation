<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# LeetCode MCP Server Project

This is a Model Context Protocol (MCP) server written in Go that integrates with VS Code to fetch daily LeetCode challenges.

## Project Overview
- **Language**: Go
- **Purpose**: MCP server for retrieving LeetCode daily challenges
- **Integration**: VS Code via MCP protocol

## Development Guidelines
- Follow Go best practices and idiomatic code patterns
- Use proper error handling with descriptive error messages
- Maintain JSON-RPC 2.0 protocol compliance for MCP communication
- Keep HTTP requests efficient with proper timeouts
- Format LeetCode problem content for better readability

## Key Components
- MCP protocol handlers (initialize, tools/list, tools/call)
- LeetCode GraphQL API integration
- HTML content cleaning and formatting
- JSON-RPC communication over stdin/stdout

## Testing
- Test MCP protocol compliance
- Verify LeetCode API integration
- Check error handling scenarios

You can find more info and examples at https://modelcontextprotocol.io/llms-full.txt
