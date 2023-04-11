package temporalstriker

import (
	"context"
	"crypto/tls"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"

	"github.com/abhishekhugetech/temporalstriker/bench"
	"github.com/abhishekhugetech/temporalstriker/types"
	"github.com/abhishekhugetech/temporalstriker/utils"
)

func Start(config types.MaruConfig) {
	// Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger.Info("Zap logger created")

	// Setup tls config from environment variables
	tlsConfig, err := utils.GetTLSConfig(config.TemporalHostPort, logger)
	if err != nil {
		logger.Fatal("failed to build tls config", zap.Error(err))
	}

	// Set sticky cache for workers
	stickyCacheSize := config.StickyCacheSize
	worker.SetStickyWorkflowCacheSize(stickyCacheSize)

	// Start the bench worker
	startBenchWorker(config, logger, tlsConfig)

	select {}
}

func startBenchWorker(
	config types.MaruConfig,
	logger *zap.Logger,
	tlsConfig *tls.Config,
) {
	if !config.SkipNamespaceCreation {
		utils.CreateNamespaceIfNeeded(logger, config.Namespace, config.TemporalHostPort, tlsConfig)
	}

	serviceClient := config.Client

	// Starting bench worker
	workerOptions := worker.Options{
		BackgroundActivityContext:               context.Background(),
		MaxConcurrentWorkflowTaskPollers:        config.MaxWorkflowTaskPollers,
		MaxConcurrentActivityTaskPollers:        8 * config.MaxWorkflowTaskPollers,
		MaxConcurrentWorkflowTaskExecutionSize:  256,
		MaxConcurrentLocalActivityExecutionSize: 256,
		MaxConcurrentActivityExecutionSize:      256,
	}
	worker := worker.New(serviceClient, "temporal-bench", workerOptions)
	worker.RegisterWorkflowWithOptions(bench.Workflow, workflow.RegisterOptions{Name: "bench-workflow"})
	worker.RegisterActivityWithOptions(bench.NewActivities(serviceClient), activity.RegisterOptions{Name: "bench-"})
	err := worker.Start()
	if err != nil {
		logger.Fatal("Unable to start bench worker", zap.Error(err))
	}
}
