package service

import "testing"

func TestDefaultHuosanyunCatalogContainsAllowedRealIDsAndAliases(t *testing.T) {
	catalog := DefaultHuosanyunCatalog()
	seen := make(map[string]HuosanyunModelSpec, len(catalog))
	for _, spec := range catalog {
		seen[spec.Model] = spec
		if IsExcludedHuosanyunModel(spec.Model) {
			t.Fatalf("catalog must not contain excluded model %s", spec.Model)
		}
	}

	for _, model := range []string{"deepseek-v4-flash", "deepseek-v4-pro", "glm-4.7", "glm-5", "glm-5.1", "MiniMax-M2.5", "qwen-plus", "qwen3.6-plus", "kimi-k2.5", "kimi-k2.6"} {
		if _, ok := seen[model]; !ok {
			t.Fatalf("expected catalog to contain %s", model)
		}
	}
	for _, model := range []string{"deepseek-v3.2", "deepseek-v3", "kling-v2-1"} {
		if _, ok := seen[model]; ok {
			t.Fatalf("catalog must not contain excluded model %s", model)
		}
	}

	mapping := HuosanyunAliasMapping()[PlatformOpenAI]
	tests := map[string]string{
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
	if got := mapping["DeepSeek-V3"]; got != "" {
		t.Fatalf("DeepSeek-V3 should not be mapped because deepseek-v3.2 is excluded, got %q", got)
	}
}

func TestMergeHuosanyunCatalogPreservesCustomPricingWithoutOverwrite(t *testing.T) {
	customPrice := 0.123
	existing := []ChannelModelPricing{
		{
			Platform:    PlatformOpenAI,
			Models:      []string{"deepseek-v4-pro"},
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
			if model == "deepseek-v4-flash" {
				return
			}
		}
	}
	t.Fatalf("expected missing Huosanyun catalog entries to be appended")
}

func TestMergeHuosanyunCatalogOverwriteRemovesExcludedModels(t *testing.T) {
	existing := []ChannelModelPricing{
		{Platform: PlatformOpenAI, Models: []string{"deepseek-v3.2"}, BillingMode: BillingModeToken},
		{Platform: PlatformOpenAI, Models: []string{"kling-v2-1"}, BillingMode: BillingModePerRequest},
	}
	merged := MergeHuosanyunCatalog(existing, true)
	for _, pricing := range merged {
		for _, model := range pricing.Models {
			if IsExcludedHuosanyunModel(model) {
				t.Fatalf("excluded model %s survived overwrite merge", model)
			}
		}
	}
}

func TestMergeHuosanyunAliasesHonorsOverwrite(t *testing.T) {
	existing := map[string]map[string]string{
		PlatformOpenAI: {
			"DeepSeek-Pro": "custom-upstream",
			"DeepSeek-V3":  "custom-excluded",
		},
	}

	kept := MergeHuosanyunAliases(existing, false)
	if got := kept[PlatformOpenAI]["DeepSeek-Pro"]; got != "custom-upstream" {
		t.Fatalf("alias overwrite disabled: got %q", got)
	}

	overwritten := MergeHuosanyunAliases(existing, true)
	if got := overwritten[PlatformOpenAI]["DeepSeek-Pro"]; got != "deepseek-v4-pro" {
		t.Fatalf("alias overwrite enabled: got %q", got)
	}
	if got := overwritten[PlatformOpenAI]["DeepSeek-V3"]; got != "" {
		t.Fatalf("excluded alias should be removed on overwrite, got %q", got)
	}
}
