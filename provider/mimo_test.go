package provider

import (
	"testing"

	muxcfg "github.com/GrayCodeAI/mux/config"
)

func TestGetOrCreateProvider_XiaomiTokenPlanUsesMimoBase(t *testing.T) {
	t.Setenv("MUX_CONFIG_DIR", t.TempDir())
	if err := muxcfg.SaveProviderConfig(&muxcfg.ProviderConfig{
		XiaomiMimoTokenPlanRegion: "sgp",
	}, ""); err != nil {
		t.Fatalf("SaveProviderConfig: %v", err)
	}

	c := Client(&MuxConfig{Provider: "xiaomi_mimo_token_plan", APIKey: "tp-test-key"})
	p, err := c.getOrCreateProvider("xiaomi_mimo_token_plan")
	if err != nil {
		t.Fatalf("getOrCreateProvider: %v", err)
	}
	mimo, ok := p.(*MiMoClient)
	if !ok {
		t.Fatalf("provider type = %T, want *MiMoClient", p)
	}
	if mimo.ProviderID() != "xiaomi_mimo_token_plan" {
		t.Fatalf("providerID = %q, want xiaomi_mimo_token_plan", mimo.ProviderID())
	}
}
