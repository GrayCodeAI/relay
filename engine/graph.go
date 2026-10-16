package engine

import "github.com/GrayCodeAI/relay/operationsgraph"

// OperationsGraphInput is the host-facing input for Relay's portable
// operations graph projection.
type OperationsGraphInput = operationsgraph.Input

// OperationsGraphExport is Relay's portable operations graph projection.
type OperationsGraphExport = operationsgraph.Export

// BuildOperationsGraph projects route and normalized usage telemetry without
// exposing provider, model, request, or generated content values.
func BuildOperationsGraph(input OperationsGraphInput) (*OperationsGraphExport, error) {
	return operationsgraph.Build(input)
}
