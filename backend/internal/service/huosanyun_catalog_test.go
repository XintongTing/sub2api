package service

import "testing"

func TestDefaultHuosanyunCatalogContainsRealIDsAndAliases(t *testing.T) {
	catalog := DefaultHuosanyunCatalog()
	seen := make(map[string]HuosanyunModelSpec, len(catalog))
	for _, spec := range catalog {
		seen[spec.Model] = spec
	}

	for _, model := range []string{"deepseek-v3.2", "deepseek-v4-pro", "glm-4.7", "glm-5", "MiniMax-M2.5", "qwen-plus", "kling-v2-1"} {
		if _, ok := seen[model]; !ok {
			t.Fatalf("expected catalog to contain %s", model)
		}
	}

	mapping := HuosanyunAliasMapping()[PlatformOpenAI]
	tests := map[string]string{
		"DeepSeek-V3":  "deepseek-v3.2",
		"DeepSeek-Pro": "deepseek-v4-pro",
		"GLM4-Air":     "glm-4.7",
		"GLM4-Plus":    "glm-5",
		"Qwen3-Turbo":  "qwen-plus",
	}
	for alias, want := range tests {
		if got := mapping[alias]; got != want {
			t.Fatalf("alias %s mapped to %q, want %q", alias, got, want)
		}
	}
}

func TestMergeHuosanyunCatalogPreservesCustomPricingWithoutOverwrite(t *testing.T) {
	customPrice := 0.123
	existing := []ChannelModelPricing{
		{
			Platform:    PlatformOpenAI,
			Models:      []string{"deepseek-v3.2"},
			BillingMode: BillingModeToken,
			InputPrice:  &customPrice,
		},
		{
			Platform:    PlatformOpenAI,
			Models:      []string{"custom-model"},
			BillingMode: BillingModeToken,
			InputPrice:  &customPrice,
		},
	}

	merged := MergeHuosanyunCatalog(existing, false)
	if merged[0].InputPrice == nil || *merged[0].InputPrice != customPrice {
		t.Fatalf("existing pricing was modified without overwrite")
	}
	if merged[1].Models[0] != "custom-model" {
		t.Fatalf("custom model should be preserved")
	}
	for _, pricing := range merged {
		for _, model := range pricing.Models {
			if model == "deepseek-v4-pro" {
				return
			}
		}
	}
	t.Fatalf("expected missing Huosanyun catalog entries to be appended")
}

func TestMergeHuosanyunAliasesHonorsOverwrite(t *testing.T) {
	existing := map[string]map[string]string{
		PlatformOpenAI: {
			"DeepSeek-V3": "custom-upstream",
		},
	}

	kept := MergeHuosanyunAliases(existing, false)
	if got := kept[PlatformOpenAI]["DeepSeek-V3"]; got != "custom-upstream" {
		t.Fatalf("alias overwrite disabled: got %q", got)
	}

	overwritten := MergeHuosanyunAliases(existing, true)
	if got := overwritten[PlatformOpenAI]["DeepSeek-V3"]; got != "deepseek-v3.2" {
		t.Fatalf("alias overwrite enabled: got %q", got)
	}
}
