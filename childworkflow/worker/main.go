package main

import (
	"os"
	"os/signal"

	"go.temporal.io/temporal/client"
	"go.temporal.io/temporal/worker"
	"go.uber.org/zap"

	"github.com/temporalio/temporal-go-samples/childworkflow"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		logger.Fatal("Unable to create client", zap.Error(err))
	}

	w := worker.New(c, "child-workflow-task-list", worker.Options{
		Logger: logger,
	})

	w.RegisterWorkflow(childworkflow.SampleParentWorkflow)
	w.RegisterWorkflow(childworkflow.SampleChildWorkflow)

	err = w.Start()
	if err != nil {
		logger.Fatal("Unable to start worker", zap.Error(err))
	}

	// The workers are supposed to be long running process that should not exit.
	waitCtrlC()
	// Stop worker, close connection, clean up resources.
	w.Stop()
	_ = c.CloseConnection()
}

func waitCtrlC() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
