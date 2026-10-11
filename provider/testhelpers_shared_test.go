package provider

// userMsg builds a single-user-message conversation. Shared by tests that
// previously relied on the helper defined in the (since-moved) embedding
// cache tests.
func userMsg(s string) []MuxMessage {
	return []MuxMessage{{Role: "user", Content: s}}
}
