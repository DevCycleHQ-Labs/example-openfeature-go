package main

import (
	devcycle "github.com/devcyclehq/go-server-sdk/v2"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"log"
	"os"
	"time"
)

var devcycleClient *devcycle.Client
var openFeatureClient *openfeature.Client

func initalizeDevCycle() *devcycle.Client {
	sdkKey := os.Getenv("DEVCYCLE_SERVER_SDK_KEY")

	if len(sdkKey) == 0 {
		log.Fatalf("Add your DEVCYCLE_SERVER_SDK_KEY to the .env file")
	}

	options := devcycle.Options{
		EnableEdgeDB:                 false,
		EnableCloudBucketing:         false,
		EventFlushIntervalMS:         5 * time.Second,
		ConfigPollingIntervalMS:      5 * time.Second,
		RequestTimeout:               30 * time.Second,
		DisableAutomaticEventLogging: false,
		DisableCustomEventLogging:    false,
	}

	client, err := devcycle.NewClient(sdkKey, &options)
	if err != nil {
		log.Fatalf("Error initializing DevCycle client: %v", err)
	}

	return client
}

func getDevCycleClient() *devcycle.Client {
	if devcycleClient == nil {
		devcycleClient = initalizeDevCycle()
	}

	return devcycleClient
}

func getOpenFeatureClient() *openfeature.Client {
	if openFeatureClient == nil {
		openFeatureClient = initalizeOpenFeature()
	}
	return openFeatureClient
}

func initalizeOpenFeature() *openfeature.Client {
	if err := openfeature.SetProvider(getDevCycleClient().OpenFeatureProvider()); err != nil {
		log.Fatalf("Failed to set DevCycle provider: %v", err)
	}
	client := openfeature.NewClient("devcycle")
	// setting service user globally instead of per request - and overriding it in the request context.
	client.SetEvaluationContext(openfeature.NewEvaluationContext("api-service", map[string]interface{}{}))
	return client
}
