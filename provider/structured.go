package provider

import (
	"context"

	"github.com/GrayCodeAI/mux/provider/core"
	"github.com/GrayCodeAI/mux/provider/media"
)

// ChatWithStructuredOutput validates structured JSON responses while retaining
// MuxClient's provider resolution and default-model behavior. The schema,
// validation, prompting, and retry implementation lives in provider/media.
func (c *MuxClient) ChatWithStructuredOutput(ctx context.Context, messages []core.MuxMessage, opts core.ChatOptions, validation core.SchemaValidation) (*core.MuxResponse, error) {
	providerName := opts.Provider
	if providerName == "" {
		providerName = c.defaultProvider
	}
	p, err := c.getOrCreateProvider(providerName)
	if err != nil {
		return nil, err
	}
	if opts.Model == "" {
		opts.Model = ResolveDefaultModel(providerName)
	}
	return media.ChatWithStructuredOutput(ctx, p, providerName, messages, opts, validation)
}
