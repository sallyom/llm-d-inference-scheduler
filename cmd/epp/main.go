/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/**
 * This file is adapted from Gateway API Inference Extension
 * Original source: https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/main/cmd/epp/main.go
 * Licensed under the Apache License, Version 2.0
 */

// Package main contains the "Endpoint Picker (EPP)" program for scheduling
// inference requests.
package main

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/gateway-api-inference-extension/cmd/epp/runner"

	"github.com/llm-d/llm-d-inference-scheduler/pkg/plugins"
)

func main() {
	ctx := ctrl.SetupSignalHandler()

	// Add startup span to verify tracing is working
	tracer := otel.GetTracerProvider().Tracer("llm-d-inference-scheduler")
	ctx, span := tracer.Start(ctx, "llm_d.epp.startup")
	span.SetAttributes(
		attribute.String("component", "llm-d-inference-scheduler"),
		attribute.String("operation", "startup"),
	)
	defer span.End()

	// Register llm-d-inference-scheduler plugins
	plugins.RegisterAllPlugins()

	// Note: GIE built-in plugins are automatically registered by the runner
	// when it processes configuration in runner.parsePluginsConfiguration()

	if err := runner.NewRunner().Run(ctx); err != nil {
		span.SetAttributes(attribute.String("operation.outcome", "error"))
		os.Exit(1)
	}
	span.SetAttributes(attribute.String("operation.outcome", "success"))
}
