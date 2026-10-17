package provider

import (
	"testing"

	relaycfg "github.com/GrayCodeAI/relay/config"
)

func TestGetOrCreateProvider_XiaomiTokenPlanUsesMimoBase(t *testing.T) {
	t.Setenv("RELAY_CONFIG_DIR", t.TempDir())
	if err := relaycfg.SaveProviderConfig(&relaycfg.ProviderConfig{
		XiaomiMimoTokenPlanRegion: "sgp",
	}, ""); err != nil {
		t.Fatalf("SaveProviderConfig: %v", err)
	}

	c := Client(&RelayConfig{Provider: "xiaomi_mimo_token_plan", APIKey: "tp-test-key"})
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
