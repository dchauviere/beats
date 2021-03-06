// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// +build mage

package main

import (
	"context"

	devtools "github.com/elastic/beats/dev-tools/mage"
)

// Build builds the Beat binary.
func Build() error {
	return devtools.Build(devtools.DefaultBuildArgs())
}

// Clean cleans all generated files and build artifacts.
func Clean() error {
	return devtools.Clean()
}

// Fields generates a fields.yml for the Beat.
func Fields() error {
	return devtools.GenerateFieldsYAML()
}

// GoTestUnit executes the Go unit tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoTestUnit(ctx context.Context) error {
	return devtools.GoTest(ctx, devtools.DefaultGoTestUnitArgs())
}

// GoTestIntegration executes the Go integration tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoTestIntegration(ctx context.Context) error {
	return devtools.GoTest(ctx, devtools.DefaultGoTestIntegrationArgs())
}
