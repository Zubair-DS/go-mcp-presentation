package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// MCP Protocol structures
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Tool definitions
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type ToolSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// LeetCode API structures
type LeetCodeProblem struct {
	QuestionID         string `json:"questionId"`
	QuestionFrontendID string `json:"questionFrontendId"`
	Title              string `json:"title"`
	TitleSlug          string `json:"titleSlug"`
	Content            string `json:"content"`
	Difficulty         string `json:"difficulty"`
	Stats              string `json:"stats"`
}

type LeetCodeDailyChallenge struct {
	Date       string          `json:"date"`
	UserStatus string          `json:"userStatus"`
	Link       string          `json:"link"`
	Question   LeetCodeProblem `json:"question"`
}

type GraphQLResponse struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion LeetCodeDailyChallenge `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}

// MCP Server
type MCPServer struct {
	client *http.Client
}

func NewMCPServer() *MCPServer {
	return &MCPServer{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *MCPServer) handleListTools(req MCPRequest) MCPResponse {
	tools := []Tool{
		{
			Name:        "get_leetcode_daily_challenge",
			Description: "Fetches today's LeetCode daily challenge problem with title, difficulty, and description",
			InputSchema: ToolSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"include_content": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to include the full problem description (default: true)",
					},
				},
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

func (s *MCPServer) getDailyChallenge(includeContent bool) (interface{}, error) {
	query := `
	{
		activeDailyCodingChallengeQuestion {
			date
			userStatus
			link
			question {
				questionId
				questionFrontendId
				title
				titleSlug
				content
				difficulty
				stats
			}
		}
	}`

	payload := map[string]interface{}{
		"query": query,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL query: %w", err)
	}

	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "LeetCode-MCP-Server/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LeetCode API returned status %d", resp.StatusCode)
	}

	var graphqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	challenge := graphqlResp.Data.ActiveDailyCodingChallengeQuestion

	// Format the response
	result := map[string]interface{}{
		"date":       challenge.Date,
		"title":      challenge.Question.Title,
		"difficulty": challenge.Question.Difficulty,
		"link":       fmt.Sprintf("https://leetcode.com%s", challenge.Link),
		"problem_id": challenge.Question.QuestionFrontendID,
	}

	if includeContent && challenge.Question.Content != "" {
		// Clean up the HTML content for better readability
		content := strings.ReplaceAll(challenge.Question.Content, "</p>", "\n\n")
		content = strings.ReplaceAll(content, "<p>", "")
		content = strings.ReplaceAll(content, "<strong>", "**")
		content = strings.ReplaceAll(content, "</strong>", "**")
		content = strings.ReplaceAll(content, "<code>", "`")
		content = strings.ReplaceAll(content, "</code>", "`")
		content = strings.ReplaceAll(content, "<pre>", "\n```\n")
		content = strings.ReplaceAll(content, "</pre>", "\n```\n")
		content = strings.ReplaceAll(content, "<ul>", "\n")
		content = strings.ReplaceAll(content, "</ul>", "\n")
		content = strings.ReplaceAll(content, "<li>", "• ")
		content = strings.ReplaceAll(content, "</li>", "\n")
		content = strings.ReplaceAll(content, "&lt;", "<")
		content = strings.ReplaceAll(content, "&gt;", ">")
		content = strings.ReplaceAll(content, "&amp;", "&")
		content = strings.ReplaceAll(content, "&quot;", "\"")

		result["description"] = strings.TrimSpace(content)
	}

	return result, nil
}

func (s *MCPServer) handleToolCall(req MCPRequest) MCPResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	name, ok := params["name"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: MCPError{
				Code:    -32602,
				Message: "Missing tool name",
			},
		}
	}

	switch name {
	case "get_leetcode_daily_challenge":
		includeContent := true
		if args, ok := params["arguments"].(map[string]interface{}); ok {
			if ic, exists := args["include_content"].(bool); exists {
				includeContent = ic
			}
		}

		result, err := s.getDailyChallenge(includeContent)
		if err != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: MCPError{
					Code:    -32603,
					Message: fmt.Sprintf("Failed to fetch daily challenge: %v", err),
				},
			}
		}

		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf("LeetCode Daily Challenge retrieved successfully:\n\n%s", formatResult(result)),
					},
				},
			},
		}

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: MCPError{
				Code:    -32601,
				Message: "Tool not found",
			},
		}
	}
}

func formatResult(result interface{}) string {
	data, ok := result.(map[string]interface{})
	if !ok {
		return "Unable to format result"
	}

	var output strings.Builder
	if title, ok := data["title"].(string); ok {
		output.WriteString(fmt.Sprintf("**%s**\n", title))
	}
	if difficulty, ok := data["difficulty"].(string); ok {
		output.WriteString(fmt.Sprintf("**Difficulty:** %s\n", difficulty))
	}
	if problemId, ok := data["problem_id"].(string); ok {
		output.WriteString(fmt.Sprintf("**Problem ID:** %s\n", problemId))
	}
	if link, ok := data["link"].(string); ok {
		output.WriteString(fmt.Sprintf("**Link:** %s\n", link))
	}
	if date, ok := data["date"].(string); ok {
		output.WriteString(fmt.Sprintf("**Date:** %s\n", date))
	}
	if description, ok := data["description"].(string); ok {
		output.WriteString(fmt.Sprintf("\n**Problem Description:**\n%s\n", description))
	}

	return output.String()
}

func (s *MCPServer) handleInitialize(req MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{
					"listChanged": false,
				},
			},
			"serverInfo": map[string]interface{}{
				"name":    "leetcode-mcp-server",
				"version": "1.0.0",
			},
		},
	}
}

func (s *MCPServer) handleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleListTools(req)
	case "tools/call":
		return s.handleToolCall(req)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func main() {
	server := NewMCPServer()
	scanner := bufio.NewScanner(os.Stdin)

	log.Println("LeetCode MCP Server started")

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("Failed to parse request: %v", err)
			continue
		}

		response := server.handleRequest(req)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			continue
		}

		fmt.Println(string(responseJSON))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from stdin: %v", err)
	}
}
