package temporalstriker

import (
	"time"

	"go.temporal.io/sdk/client"
)

type MaruConfig struct {
	Client                 client.Client // temporal client object
	Namespace              string        // namespace in which the entire test would be done
	TaskQueue              string        // task queue to use
	TemporalHostPort       string        // to get the temporal tls certificate (host:port) (client.DefaultHostPort)
	SkipNamespaceCreation  bool          // weather to create the namespace or not
	StickyCacheSize        int           // Stricky cache size for worker (2048) (def:10K)
	MaxWorkflowTaskPollers int           // MaxConcurrentWorkflowTaskPollers of worker options
	NamespaceRetention     time.Duration // NamespaceRetention is the time duration after which the load testing workflows will be removed
}
