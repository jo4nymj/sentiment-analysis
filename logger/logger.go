package logger

import (
	"context"
	"log"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
)

func init() {
	loggingClient, err := stackdriverClient()
	if err != nil {
		log.Fatal("Failed creating the logger client")
	}
	Logger = loggingClient.Logger("sentiments-log")
}

var Logger *logging.Logger

func Print(message string, args ...interface{}) {
	Logger.StandardLogger(logging.Info).Println(message, args)
}

func Error(message string, args ...interface{}) {
	Logger.StandardLogger(logging.Error).Println(message, args)
}

func stackdriverClient() (client *logging.Client, err error) {
	var projectID string
	if projectID, err = metadata.ProjectID(); err == nil {
		client, err = logging.NewClient(context.Background(), projectID)
		return client, err
	}
	client, err = logging.NewClient(context.Background(), "sentiments-analysis-263717")
	return client, err
}
