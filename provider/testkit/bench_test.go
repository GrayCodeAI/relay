package testkit

import (
	"context"
	"testing"

	"github.com/GrayCodeAI/mux/provider/core"
)

func BenchmarkMockProviderScriptedChat(b *testing.B) {
	mock := NewMockProvider(MockModeEcho)
	mock.Script = []MockTurn{
		{Content: "partial response", FinishReason: "tool_use", ToolCalls: []core.ToolCall{{ID: "tc1", Name: "search", Arguments: map[string]any{"q": "x"}}}},
		{Content: "done", FinishReason: "end_turn"},
	}
	msgs := []core.MuxMessage{{Role: "user", Content: "run"}}
	opts := core.ChatOptions{Model: "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mock.Reset()
		mock.Loop = true
		if _, err := mock.Chat(context.Background(), msgs, opts); err != nil {
			b.Fatal(err)
		}
	}
}
