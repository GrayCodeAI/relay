package provider

import (
	"context"

	"github.com/GrayCodeAI/relay/provider/core"
	"github.com/GrayCodeAI/relay/provider/media"
)

// ChatWithStructuredOutput validates structured JSON responses while retaining
// RelayClient's provider resolution and default-model behavior. The schema,
// validation, prompting, and retry implementation lives in provider/media.
func (c *RelayClient) ChatWithStructuredOutput(ctx context.Context, messages []core.RelayMessage, opts core.ChatOptions, validation core.SchemaValidation) (*core.RelayResponse, error) {
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
