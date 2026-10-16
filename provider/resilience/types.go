package resilience

import (
	"context"

	"github.com/GrayCodeAI/relay/llm"
	"github.com/GrayCodeAI/relay/provider/core"
)

// Resilience decorators operate on the provider-neutral core contract. These
// aliases keep the feature implementation readable without creating a second
// request/response model.
type (
	Provider           = core.Provider
	RelayMessage       = core.RelayMessage
	RelayResponse      = core.RelayResponse
	RelayStreamEvent   = core.RelayStreamEvent
	RelayUsage         = core.RelayUsage
	StreamResult       = core.StreamResult
	ChatOptions        = core.ChatOptions
	Guardrails         = core.Guardrails
	GuardrailViolation = core.GuardrailViolation
	GuardrailError     = core.GuardrailError
	GuardrailRule      = core.GuardrailRule
	GuardrailType      = core.GuardrailType
	GuardrailAction    = core.GuardrailAction
	GuardrailSeverity  = core.GuardrailSeverity
)

const (
	GuardrailPII             = core.GuardrailPII
	GuardrailPromptInjection = core.GuardrailPromptInjection
	GuardrailHarmfulContent  = core.GuardrailHarmfulContent
	GuardrailSecretLeak      = core.GuardrailSecretLeak
	GuardrailCustom          = core.GuardrailCustom
	GuardrailBlock           = core.GuardrailBlock
	GuardrailRedact          = core.GuardrailRedact
	GuardrailWarn            = core.GuardrailWarn
)

func NewGuardrails(rules ...GuardrailRule) *Guardrails { return core.NewGuardrails(rules...) }

func NewStreamResult(events <-chan RelayStreamEvent, requestID string, cancel context.CancelFunc) *StreamResult {
	return llm.NewStreamResult(events, requestID, cancel)
}
