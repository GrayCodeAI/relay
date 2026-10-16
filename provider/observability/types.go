package observability

import (
	"context"

	"github.com/GrayCodeAI/relay/llm"
	"github.com/GrayCodeAI/relay/provider/core"
)

type (
	Provider         = core.Provider
	RelayMessage     = core.RelayMessage
	RelayResponse    = core.RelayResponse
	RelayStreamEvent = core.RelayStreamEvent
	RelayUsage       = core.RelayUsage
	ToolCall         = core.ToolCall
	ChatOptions      = core.ChatOptions
	StreamResult     = core.StreamResult
)

func NewStreamResult(events <-chan RelayStreamEvent, cancel context.CancelFunc) *StreamResult {
	return llm.NewStreamResult(events, "", cancel)
}

func NewStreamResultWithRequestID(events <-chan RelayStreamEvent, requestID string, cancel context.CancelFunc) *StreamResult {
	return llm.NewStreamResult(events, requestID, cancel)
}
