package types

import (
	"go.temporal.io/sdk/client"
)

type MaruConfig struct {
	Client                 client.Client // temporal client object
	Namespace              string        // namespace in which the entire test would be done
	TemporalHostPort       string        // to get the temporal tls certificate (host:port) (client.DefaultHostPort)
	SkipNamespaceCreation  bool          // weather to create the namespace or not
	StickyCacheSize        int           // Stricky cache size for worker (2048) (def:10K)
	PrometheusPort         int           // Port on which the prometheus metrics will be shared
	MaxWorkflowTaskPollers int           // MaxConcurrentWorkflowTaskPollers of worker options
}
