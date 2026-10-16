package provider

import (
	"context"
	"fmt"

	"github.com/GrayCodeAI/relay/provider/adapters"
	"github.com/GrayCodeAI/relay/provider/embeddings"
)

// Compile-time check that *adapters.OpenAIClient implements embeddings.Embedder.
var _ embeddings.Embedder = (*adapters.OpenAIClient)(nil)

// CreateEmbedding sends an embedding request to the specified (or default) provider.
func (c *RelayClient) CreateEmbedding(ctx context.Context, req embeddings.EmbeddingRequest, provider string) (*embeddings.EmbeddingResponse, error) {
	if provider == "" {
		provider = c.defaultProvider
	}
	p, err := c.getOrCreateProvider(provider)
	if err != nil {
		return nil, err
	}
	embedder, ok := p.(embeddings.Embedder)
	if !ok {
		return nil, fmt.Errorf("relay: provider %s does not support embeddings", provider)
	}
	return embedder.CreateEmbedding(ctx, req)
}
