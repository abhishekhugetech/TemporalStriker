package sample

import (
	"context"
	"log"

	"github.com/abhishekhugetech/temporalstriker"
	"github.com/abhishekhugetech/temporalstriker/types"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func StartMaruClient() {
	namespace := "benchtest"
	hostPort := "127.0.0.1:7233"
	prometheusPort := 9090
	numDecisionPollers := 10

	// Creating a new temporal client
	serviceClient, err := client.Dial(client.Options{
		Namespace: namespace,
		HostPort:  hostPort,
		ConnectionOptions: client.ConnectionOptions{
			TLS: nil,
		},
	})
	if err != nil {
		log.Fatalln("failed to connect temporal", err)
	}

	// need to create dummy workflow
	maruConfig := types.MaruConfig{
		Client:                 serviceClient,
		Namespace:              namespace,
		TemporalHostPort:       hostPort,
		SkipNamespaceCreation:  false,
		TaskQueue:              "temporal-bench",
		StickyCacheSize:        10000,
		PrometheusPort:         prometheusPort,
		MaxWorkflowTaskPollers: numDecisionPollers,
	}

	// Register our activities and workflows that we want to test using striker
	workerOptions := worker.Options{
		BackgroundActivityContext:               context.Background(),
		MaxConcurrentWorkflowTaskPollers:        numDecisionPollers,
		MaxConcurrentActivityTaskPollers:        8 * numDecisionPollers,
		MaxConcurrentWorkflowTaskExecutionSize:  256,
		MaxConcurrentLocalActivityExecutionSize: 256,
		MaxConcurrentActivityExecutionSize:      256,
	}
	w := worker.New(serviceClient, "taskQueue", workerOptions)
	w.RegisterWorkflowWithOptions(DummyWorkflow, workflow.RegisterOptions{Name: "dummy-workflow"})
	w.RegisterActivityWithOptions(DummyActivity, activity.RegisterOptions{Name: "dummy-activity"})
	err = w.Start()
	if err != nil {
		log.Fatalln("Unable to start default worker", err)
	}

	// Starting the striker
	temporalstriker.Start(maruConfig)

}
