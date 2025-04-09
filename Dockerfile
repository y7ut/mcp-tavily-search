FROM golang:1.23 AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=1 go build -ldflags '-extldflags "-static"' -o mcp-tavily-search .

FROM alpine AS production

RUN apk add --no-cache tzdata libc6-compat

WORKDIR /app

ENV TRVILY_API_KEY=api_key

ENV TZ=Asia/Shanghai

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=build /app/mcp-tavily-search /app/

ENTRYPOINT [ "./mcp-tavily-search" ]

# docker run --rm -i --env TRVILY_API_KEY=xxx mcp-tavily-search:latest
# or
# docker run --rm -i mcp-tavily-search:latest run xxx
CMD [ "run" ]

