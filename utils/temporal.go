package utils

import (
	"context"
	"crypto/tls"
	"time"

	"go.temporal.io/api/serviceerror"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

func CreateNamespaceIfNeeded(logger *zap.Logger, namespace string, hostPort string, tlsConfig *tls.Config) {
	logger.Info("Creating namespace", zap.String("namespace", namespace), zap.String("hostPort", hostPort))

	createNamespace := func() error {
		namespaceClient, err := client.NewNamespaceClient(client.Options{
			HostPort: hostPort,
			ConnectionOptions: client.ConnectionOptions{
				TLS: tlsConfig,
			},
		})
		if err != nil {
			logger.Error("failed to create Namespace Client", zap.Error(err))
			return err
		}

		defer namespaceClient.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		retention := 10 * time.Hour * 24
		err = namespaceClient.Register(ctx, &workflowservice.RegisterNamespaceRequest{
			Namespace:                        namespace,
			WorkflowExecutionRetentionPeriod: &retention,
		})

		if err == nil {
			logger.Info("Namespace created")
			return nil
		}

		if _, ok := err.(*serviceerror.NamespaceAlreadyExists); ok {
			logger.Info("Namespace already exists")
			return nil
		}

		return err
	}

	for {
		err := createNamespace()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
}
