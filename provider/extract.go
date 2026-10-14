package provider

import (
	"context"

	"github.com/GrayCodeAI/mux/provider/extraction"
)

// ExtractRelationships delegates relationship extraction to the extraction
// feature while retaining MuxClient's public convenience method.
func (c *MuxClient) ExtractRelationships(ctx context.Context, text string, opts extraction.ExtractOptions) ([]extraction.Relationship, error) {
	return extraction.ExtractRelationships(ctx, c, text, opts)
}

var _ extraction.StructuredChatter = (*MuxClient)(nil)
