package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Run with EXPORT_DECODER_FIXTURE=1 to refresh decoder/internal/catalogtest/testdata/minimal_v1.json
func TestExportDecoderCatalogFixture(t *testing.T) {
	t.Parallel()
	if os.Getenv("EXPORT_DECODER_FIXTURE") != "1" {
		t.Skip("set EXPORT_DECODER_FIXTURE=1 to export") // TODO: https://github.com/GrayCodeAI/relay/issues/30
	}
	c := SeedCatalog()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join("..", "..", "decoder", "internal", "catalogtest", "testdata", "minimal_v1.json")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(out, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
