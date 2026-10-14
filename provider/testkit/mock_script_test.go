package testkit

import (
	"context"
	"testing"
	"time"

	"github.com/GrayCodeAI/mux/llm"
	"github.com/GrayCodeAI/mux/provider/core"
)

func TestMockProviderScriptedTurns(t *testing.T) {
	mock := NewMockProvider(MockModeEcho)
	mock.Script = []MockTurn{
		{Content: "thinking", FinishReason: "tool_use", ToolCalls: []core.ToolCall{{ID: "tc1", Name: "search", Arguments: map[string]interface{}{"q": "x"}}}},
		{Content: "done", FinishReason: "end_turn"},
	}

	resp1, err := mock.Chat(context.Background(), []core.MuxMessage{{Role: "user", Content: "hello"}}, core.ChatOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if resp1.FinishReason != "tool_use" || len(resp1.ToolCalls) != 1 {
		t.Fatalf("turn1: %+v", resp1)
	}

	resp2, err := mock.Chat(context.Background(), []core.MuxMessage{{Role: "user", Content: "hello"}}, core.ChatOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if resp2.Content != "done" {
		t.Fatalf("turn2: %+v", resp2)
	}

	// Exhausted
	_, err = mock.Chat(context.Background(), []core.MuxMessage{{Role: "user", Content: "again"}}, core.ChatOptions{})
	if err == nil {
		t.Fatal("expected exhaustion error")
	}
	mock.Reset()
	// loop mode replays
	mock.Loop = true
	r, err := mock.Chat(context.Background(), []core.MuxMessage{{Role: "user", Content: "loop"}}, core.ChatOptions{})
	if err != nil || r.Content != "thinking" {
		t.Fatalf("loop: %+v %v", r, err)
	}
}

func TestMockProviderScriptedCancellation(t *testing.T) {
	mock := NewMockProvider(MockModeEcho)
	mock.Script = []MockTurn{{Content: "late", Delay: 200 * time.Millisecond}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := mock.Chat(ctx, []core.MuxMessage{{Role: "user", Content: "x"}}, core.ChatOptions{})
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}

func TestMockProviderScriptedStreamChat(t *testing.T) {
	mock := NewMockProvider(MockModeEcho)
	mock.Script = []MockTurn{{Content: "hello world", FinishReason: "end_turn"}}
	sr, err := mock.StreamChat(context.Background(), []core.MuxMessage{{Role: "user", Content: "x"}}, core.ChatOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer sr.Close()
	var text string
	for ev := range sr.Events {
		if ev.Type == "content" {
			text += ev.Content
		}
	}
	_ = text
	// ensure llm import used
	_ = llm.NewStreamResult
}
