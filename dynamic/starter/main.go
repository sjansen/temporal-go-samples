package main

import (
	"context"
	"time"

	"github.com/pborman/uuid"
	"go.temporal.io/temporal/client"
	"go.uber.org/zap"

	"github.com/temporalio/temporal-go-samples/dynamic"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		logger.Fatal("Unable to create client", zap.Error(err))
	}

	// This workflow ID can be user business logic identifier as well.
	workflowID := "dynamic_" + uuid.New()
	workflowOptions := client.StartWorkflowOptions{
		ID:                              workflowID,
		TaskList:                        "dynamic-task-list",
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, dynamic.GreetingsWorkflowName)
	if err != nil {
		logger.Error("Unable to execute workflow", zap.Error(err))
	} else {
		logger.Info("Started workflow", zap.String("WorkflowID", we.GetID()), zap.String("RunID", we.GetRunID()))
	}

	// Close connection, clean up resources.
	_ = c.CloseConnection()
}
