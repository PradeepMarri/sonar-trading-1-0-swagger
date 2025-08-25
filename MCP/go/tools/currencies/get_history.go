package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sonar-trading/mcp-server/config"
	"github.com/sonar-trading/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func Get_historyHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["from"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("from=%v", val))
		}
		if val, ok := args["to"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("to=%v", val))
		}
		if val, ok := args["date"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("date=%v", val))
		}
		if val, ok := args["amount"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("amount=%v", val))
		}
		if val, ok := args["decimal_places"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("decimal_places=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/history%s", cfg.BaseURL, queryString)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// No authentication required for this endpoint
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGet_historyTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_history",
		mcp.WithDescription("Return a historic rate for a currencies"),
		mcp.WithString("from", mcp.Required(), mcp.Description("Currency you want to convert. For example, EUR")),
		mcp.WithString("to", mcp.Required(), mcp.Description("Comma separated list of currencies codes. For example, USD")),
		mcp.WithString("date", mcp.Required(), mcp.Description("UTC date should be in the form of YYYY-MM-DD, for example, 2018-06-20. Data available from 2018-06-19 only.")),
		mcp.WithString("amount", mcp.Description("This parameter can be used to specify the amount you want to convert. If an amount is not specified then 1 is assumed.")),
		mcp.WithString("decimal_places", mcp.Description("This parameter can be used to specify the number of decimal places included in the output. If an amount is not specified then 4 is assumed.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Get_historyHandler(cfg),
	}
}
