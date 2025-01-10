package tavily

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/y7ut/mcp-tavily-search/pkg/param"
)

// TravilySearch is a singleton
var TravilySearch *TavilySearch

const (
	TopicGeneral         = "general"
	TopicNews            = "news"
	DepthBasic           = "basic"
	DepthAdvanced        = "advanced"
	DefaultDays          = 7
	TavilySearchEndpoint = "https://api.tavily.com/search"
)

type TavilySearchResquest struct {
	MaxResults int `json:"max_results"`

	IncludeImages     bool   `json:"include_images"`
	IncludeImageDesc  bool   `json:"include_image_descriptions"`
	IncludeAnswer     bool   `json:"include_answer"`
	IncludeRawContent bool   `json:"include_raw_content"`
	Query             string `json:"query"`

	ApiKey      string `json:"api_key"`
	Topic       string `json:"topic"`
	SearchDepth string `json:"search_depth"`
	Days        int    `json:"days"`

	IncludeDomains []string `json:"include_domains"`
	ExcludeDomains []string `json:"exclude_domains"`
}

type TavilySearch struct {
	ApiKey string

	IncludeDomains []string
	ExcludeDomains []string

	Debug bool
}

type TavilySearchImage struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type TavilySearchResult struct {
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	Content       string  `json:"content"`
	Score         float64 `json:"score"`
	RawContent    *string `json:"raw_content"`
	PublishedDate *string `json:"published_date"`
}

type TavilySearchResponse struct {
	Query             string               `json:"query"`
	FollowUpQuestions *string              `json:"follow_up_questions"`
	Answer            *string              `json:"answer"`
	Images            []TavilySearchImage  `json:"images"`
	Results           []TavilySearchResult `json:"results"`
	ResponseTime      float64              `json:"response_time"`
}

// Init initialize
func Init(apiKey string, debug bool, includeDomain []string, excludeDomain []string) {
	if TravilySearch == nil {
		TravilySearch = NewTavilySearch(apiKey, debug, includeDomain, excludeDomain)
	}
}

// NewTavilySearch
func NewTavilySearch(apiKey string, debug bool, includeDomain []string, excludeDomain []string) *TavilySearch {
	return &TavilySearch{
		ApiKey:         apiKey,
		Debug:          debug,
		IncludeDomains: includeDomain,
		ExcludeDomains: excludeDomain,
	}
}

// Search search from tavily with keyword and options
func Search(ctx context.Context, query string, h ...WithOptionHelper) ([]TavilySearchResult, error) {
	if TravilySearch == nil {
		return nil, fmt.Errorf("tavily search is not initialized")
	}
	return TravilySearch.Search(ctx, query, h...)
}

// Search
func (t *TavilySearch) Search(ctx context.Context, query string, h ...WithOptionHelper) ([]TavilySearchResult, error) {

	tavilyParams := NewOptionManager()
	for _, helper := range h {
		helper(tavilyParams)
	}

	tavilyReq, err := t.applyParams(*tavilyParams)
	if err != nil {
		return nil, err
	}
	tavilyReq.Query = query
	tavilyReq.ApiKey = t.ApiKey
	tavilyReq.IncludeDomains = t.IncludeDomains
	tavilyReq.ExcludeDomains = t.ExcludeDomains

	var body io.Reader
	reqbody, err := json.Marshal(tavilyReq)
	if err != nil {
		return nil, fmt.Errorf("tavily params marshal error: %v", err)
	}
	body = strings.NewReader(string(reqbody))

	if t.Debug {
		fmt.Fprintf(os.Stderr, "Tavily api input: %s\n", string(reqbody))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, TavilySearchEndpoint, body)
	if err != nil {
		return nil, fmt.Errorf("tavily api request error: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tavily API request error: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Tavily API response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	if t.Debug {
		fmt.Fprintf(os.Stderr, "Tavily API output: %s\n", string(respBody))
	}
	var tsResponse TavilySearchResponse
	if err := json.Unmarshal(respBody, &tsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Tavily API response: %v", err)
	}

	// 整理返回结果
	return tsResponse.Results, nil
}

// applyParams
// Available params:
// - debug: bool
// - limit: int
// - topic: string
// - search_depth: string
// - days: int
func (t *TavilySearch) applyParams(options OptionManager) (*TavilySearchResquest, error) {

	tavilyParams := TavilySearchResquest{}

	if err := param.Assign(&tavilyParams.MaxResults, options.GetOptionWithDefault("limit", 5)); err != nil {
		return nil, err
	}

	if err := param.Assign(&tavilyParams.Topic, options.GetOptionWithDefault("topic", TopicGeneral)); err != nil {
		return nil, err
	}
	if tavilyParams.Topic != TopicGeneral && tavilyParams.Topic != TopicNews {
		return nil, fmt.Errorf("tavily topic error: %s is not a valid topic", tavilyParams.Topic)
	}

	if err := param.Assign(&tavilyParams.SearchDepth, options.GetOptionWithDefault("search_depth", DepthBasic)); err != nil {
		return nil, err
	}
	if tavilyParams.SearchDepth != DepthBasic && tavilyParams.SearchDepth != DepthAdvanced {
		return nil, fmt.Errorf("tavily search depth error: %s is not a valid search depth", tavilyParams.SearchDepth)
	}

	if err := param.Assign(&tavilyParams.Days, options.GetOptionWithDefault("days", DefaultDays)); err != nil {
		return nil, err
	}
	if tavilyParams.Days < 1 || tavilyParams.Days > 30 {
		return nil, fmt.Errorf("tavily days error: %d is not a valid days, days must between 1 and 30", tavilyParams.Days)
	}

	return &tavilyParams, nil
}
