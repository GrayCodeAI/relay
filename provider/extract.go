package provider

import (
	"context"

	"github.com/GrayCodeAI/relay/provider/extraction"
)

// ExtractRelationships delegates relationship extraction to the extraction
// feature while retaining RelayClient's public convenience method.
func (c *RelayClient) ExtractRelationships(ctx context.Context, text string, opts extraction.ExtractOptions) ([]extraction.Relationship, error) {
	return extraction.ExtractRelationships(ctx, c, text, opts)
}

var _ extraction.StructuredChatter = (*RelayClient)(nil)
