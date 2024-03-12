package main

import (
	"context"
	"fmt"
	devcycle "github.com/devcyclehq/go-server-sdk/v2"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"os"
	"time"
)

// ServiceUser is used outside of a request context, we define a service user.
// This can contian properties unique to this service, and allows you to target
// services in the same way you would target app users.
var ServiceUser = devcycle.User{
	UserId: "api-service",
}

// Log the current DevCycle variation to the console.
func logVariation() {
	time.Sleep(500 * time.Millisecond)
	os.Stdout.Write([]byte("\n"))
	renderFrame(0)
}

func renderFrame(idx int) {
	client := getDevCycleClient()
	ofClient := getOpenFeatureClient()

	features, _ := client.AllFeatures(ServiceUser)
	variationName := "Default"
	if feat, exists := features["hello-togglebot"]; exists {
		variationName = feat.VariationName
	}

	variableValue, _ := ofClient.BooleanValue(context.Background(), "togglebot-wink", false, openfeature.NewTargetlessEvaluationContext(map[string]interface{}{}))
	wink := variableValue

	stringVariable, _ := ofClient.StringValue(context.Background(), "togglebot-speed", "off", openfeature.NewTargetlessEvaluationContext(map[string]interface{}{}))
	speed := stringVariable

	spinChars := [6]rune{'◜', '◠', '◝', '◞', '◡', '◟'}

	idx = (idx + 1) % len(spinChars)
	spinner := spinChars[idx]
	if speed == "off" {
		spinner = '○'
	}

	face := "(○ ‿ ○)"
	if wink {
		face = "(- ‿ ○)"
	}

	frame := fmt.Sprintf("%c Serving variation: %s %s", spinner, variationName, face)

	color := "blue"
	if speed == "surprise" {
		color = "rainbow"
	}

	writeToConsole(frame, color)

	timeout := 100 * time.Millisecond
	if speed == "off" || speed == "slow" {
		timeout = 500 * time.Millisecond
	}
	time.Sleep(timeout)
	renderFrame(idx)
}

func addColor(text string, color string) string {
	colors := make(map[string]string)
	colors["red"] = "\033[91m"
	colors["green"] = "\033[92m"
	colors["yellow"] = "\033[93m"
	colors["blue"] = "\033[94m"
	colors["magenta"] = "\033[95m"
	colors["rainbow"] = fmt.Sprintf("\033[38;5;%dm", time.Now().UnixNano()%230)
	endChar := "\033[0m"

	if colorCode, exists := colors[color]; exists {
		return colorCode + text + endChar
	} else {
		return text
	}
}

func writeToConsole(frame string, color string) {
	frame = addColor(frame, color)
	os.Stdout.Write([]byte("\x1b[K  " + frame + "\r"))
}
