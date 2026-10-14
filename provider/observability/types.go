package observability

import (
	"context"

	"github.com/GrayCodeAI/mux/llm"
	"github.com/GrayCodeAI/mux/provider/core"
)

type (
	Provider       = core.Provider
	MuxMessage     = core.MuxMessage
	MuxResponse    = core.MuxResponse
	MuxStreamEvent = core.MuxStreamEvent
	MuxUsage       = core.MuxUsage
	ToolCall       = core.ToolCall
	ChatOptions    = core.ChatOptions
	StreamResult   = core.StreamResult
)

func NewStreamResult(events <-chan MuxStreamEvent, cancel context.CancelFunc) *StreamResult {
	return llm.NewStreamResult(events, "", cancel)
}

func NewStreamResultWithRequestID(events <-chan MuxStreamEvent, requestID string, cancel context.CancelFunc) *StreamResult {
	return llm.NewStreamResult(events, requestID, cancel)
}
