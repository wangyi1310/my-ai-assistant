package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"time"
)

func NewMCPClient(endpoint string, ctx context.Context) *client.SSEMCPClient {
	cli, err := client.NewSSEMCPClient(endpoint + "/sse")
	if err != nil {
		fmt.Println("Failed to create client: %v", err)
	}
	defer cli.Close()

	// Start the client
	if err := cli.Start(ctx); err != nil {
		fmt.Println("Failed to start client: %v", err)
	}

	// Initialize
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "say-hello-server",
		Version: "1.0.0",
	}

	result, err := cli.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	if result.ServerInfo.Name != "say-hello-server" {
		log.Fatalf(
			"Expected server name 'say-hello-server', got '%s'",
			result.ServerInfo.Name,
		)
	}
	return cli
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cli := NewMCPClient("http://localhost:8081", ctx)
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := cli.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(tools)
}
