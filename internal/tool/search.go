package tool

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/y7ut/mcp-tavily-search/internal/tavily"
	"github.com/y7ut/mcp-tavily-search/pkg/param"
)

// TavilySearchHandler is the handler for the search tool
func TavilySearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var keyword string
	if err := param.Assign(&keyword, request.Params.Arguments["keyword"]); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("keyword error: %v", err)), nil
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

	textContents := make([]mcp.Content, len(result))
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

// TavilySearchImageHandler is the handler for the search image tool, return image content
func TavilySearchImageHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var keyword string
	if err := param.Assign(&keyword, request.Params.Arguments["keyword"]); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("keyword error: %v", err)), nil
	}

	result, err := tavily.SearchImage(
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

	imgContents := make([]mcp.Content, 0)
	wg := &sync.WaitGroup{}
	for _, news := range result {
		wg.Add(1)
		go func() {
			defer wg.Done()
			content, err := downloadImage(news.URL)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				return
			}
			imgContents = append(imgContents, content)
			imgContents = append(imgContents, mcp.TextContent{
				Type: "text",
				Text: news.Description,
			})
		}()
	}

	wg.Wait()

	return &mcp.CallToolResult{
		Content: imgContents,
	}, nil
}

func downloadImage(url string) (*mcp.ImageContent, error) {
	imgBuffer := bytes.NewBuffer([]byte{})
	imgBase64Buffer := base64.NewEncoder(base64.StdEncoding, imgBuffer)
	defer imgBase64Buffer.Close()
	img, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download image error: %v", err)
	}
	mimeType := img.Header.Get("Content-Type")
	defer img.Body.Close()

	_, err = io.Copy(imgBase64Buffer, img.Body)
	if err != nil {
		return nil, fmt.Errorf("encode image error: %v", err)
	}
	fmt.Fprintf(os.Stderr, "Downloading image: %s\n", url)
	return &mcp.ImageContent{
		Type:     "image",
		MIMEType: mimeType,
		Data:     imgBuffer.String(),
	}, nil

}
