package mcp

import (
	"fmt"
	stdlog "log"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/mark3labs/mcp-go/server"
)

type ResourceVerifierServer struct {
	cred   *azidentity.DefaultAzureCredential
	logger *stdlog.Logger
	server *server.MCPServer
}

func NewResourceVerifierServer(logFile string) (*ResourceVerifierServer, error) {
	logger, err := initLogger(logFile)
	if err != nil {
		stdlog.Fatal("Failed to initialize logger:", err)
	}

	azureSubscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if azureSubscriptionId == "" {
		logger.Fatal("AZURE_SUBSCRIPTION_ID is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	s := &ResourceVerifierServer{
		cred: cred,
		server: server.NewMCPServer(
			"azure-resource-verifier",
			"0.0.1",
			server.WithResourceCapabilities(true, true),
			server.WithToolCapabilities(true),
			server.WithLogging(),
		),
		logger: stdlog.New(logger.Writer(), "azure-resource-verifier", 0),
	}

	// Register tool handlers
	s.server.AddTool(getLocations(cred, azureSubscriptionId))

	return s, nil
}

func initLogger(outPath string) (*log.Logger, error) {
	if outPath == "" {
		return log.New(), nil
	}

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetOutput(file)

	return logger, nil
}

func (s *ResourceVerifierServer) ServeStdio() error {
	return server.ServeStdio(s.server, server.WithErrorLogger(s.logger))
}
