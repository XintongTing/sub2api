package service

import (
	"strings"
	"testing"
)

func TestDefaultHuosanyunCatalogContainsAllowedRealIDsAndRejectsLegacyAliases(t *testing.T) {
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
	for _, model := range []string{"deepseek-v4-pro", "glm-4.7", "glm-5", "qwen-plus"} {
		if got := mapping[model]; got != model {
			t.Fatalf("real model %s mapped to %q", model, got)
		}
	}
	for _, alias := range []string{"DeepSeek-Pro", "GLM4-Air", "GLM4-Plus", "Qwen3-Turbo", "DeepSeek-V3"} {
		if got := mapping[alias]; got != "" {
			t.Fatalf("legacy alias %s should not be mapped, got %q", alias, got)
		}
		if !IsExcludedHuosanyunModel(alias) {
			t.Fatalf("legacy alias %s should be excluded", alias)
		}
	}
}

func TestHuosanyunSeedanceCatalogIsPublicButNotChatCallableByDefault(t *testing.T) {
	foundSeedance := false
	for _, pricing := range HuosanyunCatalogPricing() {
		if len(pricing.Models) == 0 || !strings.Contains(strings.ToLower(pricing.Models[0]), "seedance") {
			continue
		}
		foundSeedance = true
		if pricing.IsAPIEnabled() {
			t.Fatalf("seedance model %s should not be exposed through chat completions by default", pricing.Models[0])
		}
	}
	if !foundSeedance {
		t.Fatalf("expected seedance catalog entries")
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

func TestMergeHuosanyunAliasesRemovesLegacyAliases(t *testing.T) {
	existing := map[string]map[string]string{
		PlatformOpenAI: {
			"DeepSeek-Pro":   "custom-upstream",
			"DeepSeek-V3":    "custom-excluded",
			"deepseek-v4-pro": "custom-real",
		},
	}

	kept := MergeHuosanyunAliases(existing, false)
	if got := kept[PlatformOpenAI]["DeepSeek-Pro"]; got != "" {
		t.Fatalf("legacy alias should be removed even when overwrite is disabled, got %q", got)
	}
	if got := kept[PlatformOpenAI]["DeepSeek-V3"]; got != "" {
		t.Fatalf("excluded alias should be removed, got %q", got)
	}
	if got := kept[PlatformOpenAI]["deepseek-v4-pro"]; got != "custom-real" {
		t.Fatalf("real mapping should be preserved without overwrite, got %q", got)
	}

	overwritten := MergeHuosanyunAliases(existing, true)
	if got := overwritten[PlatformOpenAI]["deepseek-v4-pro"]; got != "deepseek-v4-pro" {
		t.Fatalf("real mapping should be refreshed on overwrite, got %q", got)
	}
}
