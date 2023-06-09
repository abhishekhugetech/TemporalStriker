# Temporal Striker - Maru Golang SDK

Temporal Striker is a Golang SDK based on the [`Maru`](https://github.com/temporalio/maru) project that exposes Maru's codebase as a Golang SDK. With Temporal Striker, you can build applications that leverage Maru's functionality in a more efficient and effective way.

**⚠️ WARNING**

To understand and use this Tool please go through the [`Maru`](https://github.com/temporalio/maru) repository once.


## Problem It solves

Rather than just copying your workflow and activities into Maru's codebase or copying the Maru's code to you project, You can add [`Temporal Striker`]() to your project and directly start the load testing without getting into any hastle.


## Getting Started

To get started with Temporal Striker, you need to have Golang installed on your machine. Once you have Golang installed, you can use the following command to install Temporal Striker:


#### 1. Adding TemporalStriker to your project

```bash
go get github.com/abhishekhugetech/temporalstriker
```

#### 2. Setup your temporal client

```go
namespace := "benchtest"
hostPort := "127.0.0.1:7233"
 := 9090
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
```

Once we have initialized `TemporalClient` we can easily add workers and register our workflows and activities in our workers, that will be part of our load testing.


#### Registering worker and workflow

```go
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
    log.Println("Unable to start worker", err)
}
```


#### 3. Create `MaruConfig` for initializing `TemporalStriker`

```go
c := types.MaruConfig{
		Client:                 serviceClient,
		Namespace:              namespace,
		TemporalHostPort:       hostPort,
		SkipNamespaceCreation:  false,
		StickyCacheSize:        10000,
		MaxWorkflowTaskPollers: numDecisionPollers,
}
```


#### 4. Starting Maru

Once we have created the MaruConfig we can simply call the `temporalstriker.Start(config)` method by passing the `MaruConfig` that we just created.

This will start `Maru's` bench workflow which schedules the workflows instructed in the scenarios json files.


```go
temporalstriker.Start(c)
```


#### 5. Load Testing

The load testing part is same as `Maru`, we'll just have to use the `tctl` to Load test our workflows with TemporalStriker.

```bash
tctl --namespace benchtest wf start --tq temporal-bench --wt bench-workflow --wtt 5 --et 1800 --if ./sample/dummy.json --wid 1
```


The above command will simply start the `dummy-workflow` workflow that we registered with out worker when setting up the TemporalClient.


#### 6. Releasing a new update

To release a new update we need to add a tag to the current branch and push the tag to origin.

```bash
# Adding a tag
git tag v0.0.13
# pushing a tag
git push origin v0.0.13
```

## Miscellaneous

### Adding Metrics handler

For getting the information regarding the bechmarking, Please add `MetricsHandler` to you temporal client when you setup your `TemporalClient`.

```go
func newPrometheusScope(c prometheus.Configuration) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Println("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)
	return scope
}

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

```

If the Metrics handler is attached to our temporal client, We'll be able to access `Prometheus` metrics on port 9090.



---


**For see the sample implementation of TemporalStriker, Please refer to the Sample package.**


### TODO:

- [ ] test does not completes when using UUID workflow IDs