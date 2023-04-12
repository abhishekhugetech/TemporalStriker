package sample

import (
	"log"
	"time"

	"go.temporal.io/sdk/workflow"
)

func DummyWorkflow(ctx workflow.Context) (string, error) {

	log.Println("starting Dummy workflow")

	ao := workflow.ActivityOptions{
		TaskQueue:           "taskQueue",
		StartToCloseTimeout: time.Minute * 10,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, "dummy-activity").Get(ctx, &result)
	if err != nil {
		return "failed", err
	}

	return result, nil
}
