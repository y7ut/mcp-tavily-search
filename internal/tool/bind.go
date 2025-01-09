package tool

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/y7ut/mcp-tavily-search/internal/tavily"
)

// Bind binds the search tool
func Bind(server *server.MCPServer) {
	// Add tool
	searchTool := mcp.NewTool("search_news",
		mcp.WithDescription("Get recent news from tavily by keyword"),
		mcp.WithString("keyword",
			mcp.Required(),
			mcp.Description("Keyword to search for."),
		),
		mcp.WithNumber("days",
			mcp.DefaultNumber(7),
			mcp.Description("Number of days to search, default is 7 days."),
		),
		mcp.WithNumber("limit",
			mcp.DefaultNumber(5),
			mcp.Description("Number of news to return, default is 5."),
		),
		mcp.WithString("search_depth",
			mcp.Enum(tavily.DepthAdvanced, tavily.DepthBasic),
			mcp.DefaultString(tavily.DepthBasic),
			mcp.Description("The depth of the search. It can be \"basic\" or \"advanced\". Default is \"basic\" unless specified otherwise in a given method. "),
		),
		mcp.WithString("topic",
			mcp.Enum(tavily.TopicGeneral, tavily.TopicNews),
			mcp.DefaultString(tavily.TopicNews),
			mcp.Description("The topic of the search, default is news. topic news will retrun high quality news, topic general will return unprocessed website pages."),
		),
	)
	server.AddTool(searchTool, TavilySearchHandler)
}
