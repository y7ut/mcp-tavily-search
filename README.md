# MCP TAVILY SEARCH

A Model Context Protocol (MCP) server that provide search by tavily.

## Quick start

install

```sh
go install github.com/y7ut/mcp-tavily-search@latest
```

add config to mcp config file.

```json
{
  "mcpServers": {
    "tavily": {
      "command": "mcp-tavily-search",
      "args": [
        "run",
        "tvly-*******************"
      ]
    }
  }
}
```

or debug

```sh
npx @modelcontextprotocol/inspector mcp-tavily-search run tvly-xxxxxxxxxx
```

## Tools

### search_news

| **Parameter**   | **Type**   | **Default Value** | **Description**                                                                                                                                           | **Required** |
|------------------|------------|-------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|--------------|
| `keyword`        | `string`   | N/A               | The keyword to search for.                                                                                                                                | Yes          |
| `days`           | `number`   | `7`               | Number of days to search within. Default is 7 days.                                                                                                       | No           |
| `limit`          | `number`   | `5`               | Number of news articles to return. Default is 5.                                                                                                          | No           |
| `search_depth`   | `string`   | `"basic"`         | The depth of the search. It can be `"basic"` or `"advanced"`. Default is `"basic"`.                                                                       | No           |
| `topic`          | `string`   | `"news"`          | The topic of the search. Options are `"general"` (unprocessed pages) or `"news"` (high-quality news). Default is `"news"`.                                 | No           |
