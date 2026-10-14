package testkit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/GrayCodeAI/mux/llm"
	"github.com/GrayCodeAI/mux/provider/core"
)

// MockMode controls how the mock provider responds.
type MockMode string

const (
	// MockModeEcho echoes the last user message back.
	MockModeEcho MockMode = "echo"
	// MockModeFixed returns a fixed response set via MockProvider.Response.
	MockModeFixed MockMode = "fixed"
	// MockModeToolUse returns a tool call response.
	MockModeToolUse MockMode = "tool_use"
	// MockModeError always returns an error.
	MockModeError MockMode = "error"
	// MockModeMaxTokens returns a response with stop_reason=max_tokens (for testing continuation).
	MockModeMaxTokens MockMode = "max_tokens"
)

// MockTurn is a single scripted reply for a scripted mock provider. Each call to
// Chat/StreamChat consumes the next turn, letting tests drive multi-turn
// tool-calling conversations without a network.
type MockTurn struct {
	Content      string                                                  // assistant text
	ToolCalls    []core.ToolCall                                         // optional tool calls
	FinishReason string                                                  // defaults to end_turn / tool_use
	Usage        *core.MuxUsage                                          // defaults to a small non-zero usage
	Err          error                                                   // if set, the call fails with it
	Delay        time.Duration                                           // simulate latency (cancellation aware)
	OnBefore     func(messages []core.MuxMessage, opts core.ChatOptions) // inspect inputs, late (outside lock)
}

// MockProvider is a core.Provider implementation for testing.
// It never makes real HTTP requests.
//
// Modes (Mode field) are used when Script is empty. When Script is populated it
// is consumed turn-by-turn and Mode is ignored.
type MockProvider struct {
	mu       sync.Mutex
	Mode     MockMode
	Response string // used in MockModeFixed
	ToolName string // used in MockModeToolUse
	ToolArgs map[string]interface{}
	Delay    time.Duration // simulate latency
	Calls    []MockCall    // recorded calls for assertions

	// Script drives deterministic multi-turn conversations. When non-empty,
	// Chat/StreamChat consume one turn per call instead of Mode.
	Script []MockTurn
	// Loop replays the script from the start once it is exhausted.
	Loop bool
	// scriptIndex tracks the next Script entry to consume.
	scriptIndex int
}

// MockCall records a single call to the mock provider.
type MockCall struct {
	Messages []core.MuxMessage
	Options  core.ChatOptions
}

// Compile-time check.
var _ core.Provider = (*MockProvider)(nil)

// NewMockProvider creates a mock provider with the given mode.
func NewMockProvider(mode MockMode) *MockProvider {
	return &MockProvider{Mode: mode}
}

// Name returns "mock".
func (m *MockProvider) Name() string { return "mock" }

// Ping always succeeds.
func (m *MockProvider) Ping(_ context.Context) error { return nil }

// Chat returns a mock response based on Mode.
func (m *MockProvider) Chat(ctx context.Context, messages []core.MuxMessage, opts core.ChatOptions) (*core.MuxResponse, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, MockCall{Messages: messages, Options: opts})

	var turn MockTurn
	usedScript := false
	if len(m.Script) > 0 {
		usedScript = true
		switch {
		case m.scriptIndex < len(m.Script):
			turn = m.Script[m.scriptIndex]
			m.scriptIndex++
		case m.Loop && len(m.Script) > 0:
			turn = m.Script[m.scriptIndex%len(m.Script)]
			m.scriptIndex++
		default:
			m.mu.Unlock()
			return nil, fmt.Errorf("mux: mock script exhausted after %d turns", len(m.Script))
		}
	}
	m.mu.Unlock()

	if usedScript && turn.OnBefore != nil {
		turn.OnBefore(messages, opts)
	}

	delay := m.Delay
	if usedScript && turn.Delay > 0 {
		delay = turn.Delay
	}
	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	} else {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	if usedScript {
		if turn.Err != nil {
			return nil, turn.Err
		}
		finish := turn.FinishReason
		if finish == "" {
			if len(turn.ToolCalls) > 0 {
				finish = "tool_use"
			} else {
				finish = "end_turn"
			}
		}
		usage := turn.Usage
		if usage == nil {
			usage = &core.MuxUsage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15}
		}
		return &core.MuxResponse{
			Content:      turn.Content,
			ToolCalls:    turn.ToolCalls,
			FinishReason: finish,
			Usage:        usage,
		}, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	// Re-add the call synchronously after script branch? Already recorded above.
	switch m.Mode {
	case MockModeError:
		return nil, fmt.Errorf("mux: mock error")
	case MockModeMaxTokens:
		return &core.MuxResponse{Content: "partial response", FinishReason: "max_tokens", Usage: &core.MuxUsage{PromptTokens: 10, CompletionTokens: 100, TotalTokens: 110}}, nil
	case MockModeToolUse:
		name := m.ToolName
		if name == "" {
			name = "mock_tool"
		}
		args := m.ToolArgs
		if args == nil {
			args = map[string]interface{}{"input": "test"}
		}
		return &core.MuxResponse{
			ToolCalls:    []core.ToolCall{{ID: "mock-tc-1", Name: name, Arguments: args}},
			FinishReason: "tool_use",
			Usage:        &core.MuxUsage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
		}, nil
	case MockModeFixed:
		return &core.MuxResponse{Content: m.Response, FinishReason: "end_turn", Usage: &core.MuxUsage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15}}, nil
	default: // echo
		last := ""
		for i := len(messages) - 1; i >= 0; i-- {
			if messages[i].Role == "user" {
				last = messages[i].Content
				break
			}
		}
		return &core.MuxResponse{Content: "echo: " + last, FinishReason: "end_turn", Usage: &core.MuxUsage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15}}, nil
	}
}

// StreamChat streams a mock response word by word.
func (m *MockProvider) StreamChat(ctx context.Context, messages []core.MuxMessage, opts core.ChatOptions) (*core.StreamResult, error) {
	resp, err := m.Chat(ctx, messages, opts)
	if err != nil {
		return nil, err
	}

	streamCtx, cancel := context.WithCancel(ctx)
	ch := make(chan core.MuxStreamEvent, 64)

	go func() {
		defer close(ch)
		if resp.FinishReason == "tool_use" && len(resp.ToolCalls) > 0 {
			core.Emit(streamCtx, ch, core.MuxStreamEvent{Type: "tool_call", ToolCall: &resp.ToolCalls[0]})
		} else {
			words := strings.Fields(resp.Content)
			for _, w := range words {
				if m.Delay > 0 {
					select {
					case <-time.After(m.Delay / time.Duration(len(words)+1)):
					case <-streamCtx.Done():
						return
					}
				}
				core.Emit(streamCtx, ch, core.MuxStreamEvent{Type: "content", Content: w + " "})
			}
		}
		core.Emit(streamCtx, ch, core.MuxStreamEvent{Type: "done"})
	}()

	return llm.NewStreamResult(ch, "", cancel), nil
}

// CallCount returns the number of recorded calls.
func (m *MockProvider) CallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Calls)
}

// LastCall returns the most recent recorded call, or nil.
func (m *MockProvider) LastCall() *MockCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.Calls) == 0 {
		return nil
	}
	c := m.Calls[len(m.Calls)-1]
	return &c
}

// Reset clears recorded calls and resets the scripted turn cursor.
func (m *MockProvider) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = nil
	m.scriptIndex = 0
}

// MarshalCalls returns recorded calls as JSON for debugging.
func (m *MockProvider) MarshalCalls() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, _ := json.MarshalIndent(m.Calls, "", "  ")
	return string(b)
}
