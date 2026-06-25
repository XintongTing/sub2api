package service

import "strings"

// HuosanyunModelSpec is the local sellable catalog entry imported into channel
// pricing. Prices are user-facing THB sell prices, not upstream CNY cost prices.
type HuosanyunModelSpec struct {
	Model           string
	BillingMode     BillingMode
	InputTHBPer1M   float64
	OutputTHBPer1M  float64
	CacheReadTHB1M  float64
	CacheWriteTHB1M float64
	PerRequestTHB   float64
}

// DefaultHuosanyunCatalog is the verified model set for TokenAPIFuel's Huosanyun
// OpenAI-compatible channel. It keeps upstream model IDs as the primary names;
// legacy customer-facing aliases are intentionally excluded from public and API
// model lists so users only see and call real upstream model IDs.
func DefaultHuosanyunCatalog() []HuosanyunModelSpec {
	return []HuosanyunModelSpec{
		{Model: "deepseek-v4-flash", BillingMode: BillingModeToken, InputTHBPer1M: 15, OutputTHBPer1M: 30, CacheReadTHB1M: 0.3},
		{Model: "deepseek-v4-pro", BillingMode: BillingModeToken, InputTHBPer1M: 180, OutputTHBPer1M: 360, CacheReadTHB1M: 36},
		{Model: "doubao-seedance-2-0-fast-idle-260128", BillingMode: BillingModeToken, InputTHBPer1M: 555, OutputTHBPer1M: 555},
		{Model: "doubao-seedance-2-0-idle-260128", BillingMode: BillingModeToken, InputTHBPer1M: 765, OutputTHBPer1M: 765},
		{Model: "glm-4.7", BillingMode: BillingModeToken, InputTHBPer1M: 4.05, OutputTHBPer1M: 16.5, CacheReadTHB1M: 8.22},
		{Model: "glm-5", BillingMode: BillingModeToken, InputTHBPer1M: 11.25, OutputTHBPer1M: 36},
		{Model: "glm-5.1", BillingMode: BillingModeToken, InputTHBPer1M: 90, OutputTHBPer1M: 360, CacheReadTHB1M: 18, CacheWriteTHB1M: 112.5},
		{Model: "kimi-k2.5", BillingMode: BillingModeToken, InputTHBPer1M: 8.25, OutputTHBPer1M: 41.4},
		{Model: "kimi-k2.6", BillingMode: BillingModeToken, InputTHBPer1M: 97.5, OutputTHBPer1M: 405, CacheReadTHB1M: 15},
		{Model: "MiniMax-M2.5", BillingMode: BillingModeToken},
		{Model: "qwen-plus", BillingMode: BillingModeToken},
		{Model: "qwen3.6-plus", BillingMode: BillingModeToken},
	}
}

// IsExcludedHuosanyunModel returns true for upstream models or old public aliases
// that TokenAPIFuel must not expose by default.
func IsExcludedHuosanyunModel(model string) bool {
	name := strings.ToLower(strings.TrimSpace(model))
	if name == "" {
		return false
	}
	switch name {
	case "deepseek-v3.2", "deepseek-v3", "deepseek-pro", "qwen3-turbo", "glm4-air", "glm4-plus":
		return true
	default:
		return strings.HasPrefix(name, "kling")
	}
}

// CanonicalHuosanyunModelName keeps real upstream model IDs stable. Legacy
// aliases are not canonicalized because the customer asked to hide and reject
// those public-facing names.
func CanonicalHuosanyunModelName(model string) string {
	name := strings.TrimSpace(model)
	if name == "" {
		return ""
	}
	for _, spec := range DefaultHuosanyunCatalog() {
		if strings.EqualFold(name, spec.Model) {
			return spec.Model
		}
	}
	return name
}

func HuosanyunCatalogPricing() []ChannelModelPricing {
	specs := DefaultHuosanyunCatalog()
	out := make([]ChannelModelPricing, 0, len(specs))
	for _, spec := range specs {
		if IsExcludedHuosanyunModel(spec.Model) {
			continue
		}
		p := ChannelModelPricing{
			Platform:    PlatformOpenAI,
			Models:      []string{spec.Model},
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
		if strings.Contains(strings.ToLower(spec.Model), "seedance") {
			p.APIEnabled = boolPtr(false)
		}
		out = append(out, p)
	}
	return out
}

func HuosanyunAliasMapping() map[string]map[string]string {
	mapping := map[string]map[string]string{PlatformOpenAI: map[string]string{}}
	for _, spec := range DefaultHuosanyunCatalog() {
		if IsExcludedHuosanyunModel(spec.Model) {
			continue
		}
		mapping[PlatformOpenAI][spec.Model] = spec.Model
	}
	return mapping
}

func huosanyunCatalogRemovalKeys() map[string]struct{} {
	remove := map[string]struct{}{
		"deepseek-v3.2": {},
		"deepseek-v3":   {},
		"deepseek-pro":  {},
		"qwen3-turbo":   {},
		"glm4-air":      {},
		"glm4-plus":     {},
	}
	for _, p := range HuosanyunCatalogPricing() {
		for _, model := range p.Models {
			remove[strings.ToLower(strings.TrimSpace(model))] = struct{}{}
		}
	}
	return remove
}

func MergeHuosanyunCatalog(pricing []ChannelModelPricing, overwrite bool) []ChannelModelPricing {
	catalog := HuosanyunCatalogPricing()
	if overwrite {
		remove := huosanyunCatalogRemovalKeys()
		kept := make([]ChannelModelPricing, 0, len(pricing))
		for _, p := range pricing {
			shouldRemove := false
			for _, model := range p.Models {
				modelKey := strings.ToLower(strings.TrimSpace(model))
				if _, ok := remove[modelKey]; ok || IsExcludedHuosanyunModel(modelKey) {
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
			modelKey := strings.ToLower(strings.TrimSpace(model))
			if IsExcludedHuosanyunModel(modelKey) {
				continue
			}
			existing[modelKey] = struct{}{}
		}
	}
	out := make([]ChannelModelPricing, 0, len(pricing)+len(catalog))
	for _, p := range pricing {
		keep := make([]string, 0, len(p.Models))
		for _, model := range p.Models {
			if !IsExcludedHuosanyunModel(model) {
				keep = append(keep, model)
			}
		}
		if len(keep) == 0 {
			continue
		}
		p.Models = keep
		out = append(out, p)
	}
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
	for src, dst := range out[PlatformOpenAI] {
		if IsExcludedHuosanyunModel(src) || IsExcludedHuosanyunModel(dst) {
			delete(out[PlatformOpenAI], src)
		}
	}
	return out
}

func floatPtr(v float64) *float64 {
	return &v
}
