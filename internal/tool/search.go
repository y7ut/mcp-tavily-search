package tool

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/y7ut/mcp-tavily-search/internal/tavily"
	"github.com/y7ut/mcp-tavily-search/pkg/param"
)

// TavilySearchHandler is the handler for the search tool
func TavilySearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var keyword string
	if err := param.Assign(&keyword, request.Params.Arguments["keyword"]); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err := tavily.Search(
		ctx,
		keyword,
		tavily.WithOption("topic", request.Params.Arguments["topic"]),
		tavily.WithOption("days", request.Params.Arguments["days"]),
		tavily.WithOption("limit", request.Params.Arguments["limit"]),
		tavily.WithOption("search_depth", request.Params.Arguments["search_depth"]),
	)

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if len(result) == 0 {
		return mcp.NewToolResultError(fmt.Sprintf("no news found for keyword: %s", keyword)), nil
	}
	textContents := make([]any, len(result))
	for i, news := range result {
		textContents[i] = mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("《%s》: %s\n %s", news.Title, news.URL, news.Content),
		}
	}

	return &mcp.CallToolResult{
		Content: textContents,
	}, nil
}
