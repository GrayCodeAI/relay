package cache

import "github.com/GrayCodeAI/mux/provider/core"

// Cache decorators use the shared provider contract; cache has no dependency
// on the provider composition root.
type (
	Provider     = core.Provider
	MuxMessage   = core.MuxMessage
	MuxResponse  = core.MuxResponse
	ChatOptions  = core.ChatOptions
	StreamResult = core.StreamResult
)
