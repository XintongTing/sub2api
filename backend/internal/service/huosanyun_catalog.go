package service

import "strings"

// HuosanyunModelSpec is the local sellable catalog entry imported into channel
// pricing. Prices are user-facing THB sell prices, not upstream CNY cost prices.
type HuosanyunModelSpec struct {
	Model           string
	Aliases         []string
	BillingMode     BillingMode
	InputTHBPer1M   float64
	OutputTHBPer1M  float64
	CacheReadTHB1M  float64
	CacheWriteTHB1M float64
	PerRequestTHB   float64
}

// DefaultHuosanyunCatalog is the verified model set for TokenAPIFuel's Huosanyun
// OpenAI-compatible channel. It keeps the real upstream model IDs while retaining
// the legacy customer-facing aliases as compatibility mappings.
func DefaultHuosanyunCatalog() []HuosanyunModelSpec {
	return []HuosanyunModelSpec{
		{Model: "qwen-plus", Aliases: []string{"Qwen3-Turbo"}, BillingMode: BillingModeToken, InputTHBPer1M: 8, OutputTHBPer1M: 24},
		{Model: "deepseek-v3.2", Aliases: []string{"DeepSeek-V3"}, BillingMode: BillingModeToken, InputTHBPer1M: 30, OutputTHBPer1M: 45, CacheReadTHB1M: 6, CacheWriteTHB1M: 37.5},
		{Model: "glm-4.7", Aliases: []string{"GLM4-Air"}, BillingMode: BillingModeToken, InputTHBPer1M: 4.05, OutputTHBPer1M: 16.5, CacheReadTHB1M: 8.22},
		{Model: "glm-5", Aliases: []string{"GLM4-Plus"}, BillingMode: BillingModeToken, InputTHBPer1M: 11.25, OutputTHBPer1M: 36},
		{Model: "deepseek-v4-pro", Aliases: []string{"DeepSeek-Pro"}, BillingMode: BillingModeToken, InputTHBPer1M: 180, OutputTHBPer1M: 360, CacheReadTHB1M: 36},
	}
}

func HuosanyunCatalogPricing() []ChannelModelPricing {
	specs := DefaultHuosanyunCatalog()
	out := make([]ChannelModelPricing, 0, len(specs))
	for _, spec := range specs {
		models := append([]string{spec.Model}, spec.Aliases...)
		p := ChannelModelPricing{
			Platform:    PlatformOpenAI,
			Models:      models,
			BillingMode: spec.BillingMode,
		}
		if spec.InputTHBPer1M > 0 {
			p.InputPrice = floatPtr(spec.InputTHBPer1M / 1_000_000)
		}
		if spec.OutputTHBPer1M > 0 {
			p.OutputPrice = floatPtr(spec.OutputTHBPer1M / 1_000_000)
		}
		if spec.CacheReadTHB1M > 0 {
			p.CacheReadPrice = floatPtr(spec.CacheReadTHB1M / 1_000_000)
		}
		if spec.CacheWriteTHB1M > 0 {
			p.CacheWritePrice = floatPtr(spec.CacheWriteTHB1M / 1_000_000)
		}
		if spec.PerRequestTHB > 0 {
			p.PerRequestPrice = floatPtr(spec.PerRequestTHB)
		}
		out = append(out, p)
	}
	return out
}

func HuosanyunAliasMapping() map[string]map[string]string {
	mapping := map[string]map[string]string{PlatformOpenAI: {}}
	for _, spec := range DefaultHuosanyunCatalog() {
		mapping[PlatformOpenAI][spec.Model] = spec.Model
		for _, alias := range spec.Aliases {
			mapping[PlatformOpenAI][alias] = spec.Model
		}
	}
	return mapping
}

func MergeHuosanyunCatalog(pricing []ChannelModelPricing, overwrite bool) []ChannelModelPricing {
	catalog := HuosanyunCatalogPricing()
	if overwrite {
		remove := make(map[string]struct{})
		for _, p := range catalog {
			for _, model := range p.Models {
				remove[strings.ToLower(strings.TrimSpace(model))] = struct{}{}
			}
		}
		kept := make([]ChannelModelPricing, 0, len(pricing))
		for _, p := range pricing {
			shouldRemove := false
			for _, model := range p.Models {
				if _, ok := remove[strings.ToLower(strings.TrimSpace(model))]; ok {
					shouldRemove = true
					break
				}
			}
			if !shouldRemove {
				kept = append(kept, p)
			}
		}
		return append(kept, catalog...)
	}

	existing := make(map[string]struct{})
	for _, p := range pricing {
		for _, model := range p.Models {
			existing[strings.ToLower(strings.TrimSpace(model))] = struct{}{}
		}
	}
	out := append([]ChannelModelPricing(nil), pricing...)
	for _, p := range catalog {
		skip := false
		for _, model := range p.Models {
			if _, ok := existing[strings.ToLower(strings.TrimSpace(model))]; ok {
				skip = true
				break
			}
		}
		if !skip {
			out = append(out, p)
		}
	}
	return out
}

func MergeHuosanyunAliases(mapping map[string]map[string]string, overwrite bool) map[string]map[string]string {
	out := make(map[string]map[string]string, len(mapping)+1)
	for platform, entries := range mapping {
		out[platform] = make(map[string]string, len(entries))
		for k, v := range entries {
			out[platform][k] = v
		}
	}
	if out[PlatformOpenAI] == nil {
		out[PlatformOpenAI] = map[string]string{}
	}
	for src, dst := range HuosanyunAliasMapping()[PlatformOpenAI] {
		if _, exists := out[PlatformOpenAI][src]; !exists || overwrite {
			out[PlatformOpenAI][src] = dst
		}
	}
	return out
}

func floatPtr(v float64) *float64 {
	return &v
}
