package catalog

import "errors"

// ErrCatalogCacheRequired is returned when no valid ~/.relay/model_catalog.json exists.
// Run catalog discovery (decoder models refresh / relay catalog discover) to populate the cache.
var ErrCatalogCacheRequired = errors.New("model catalog cache required")
