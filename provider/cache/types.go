package cache

import "github.com/GrayCodeAI/relay/provider/core"

// Cache decorators use the shared provider contract; cache has no dependency
// on the provider composition root.
type (
	Provider      = core.Provider
	RelayMessage  = core.RelayMessage
	RelayResponse = core.RelayResponse
	ChatOptions   = core.ChatOptions
	StreamResult  = core.StreamResult
)
