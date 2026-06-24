package handler

import (
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AvailableChannelHandler 处理用户侧「可用渠道」查询。
//
// 用户侧接口委托 ChannelService.ListAvailable，并在返回前做三层过滤：
//  1. 行过滤：只保留状态为 Active 且与当前用户可访问分组有交集的渠道；
//  2. 分组过滤：渠道的 Groups 只保留用户可访问的那些；
//  3. 平台过滤：渠道的 SupportedModels 只保留平台在用户可见 Groups 中出现过的模型，
//     防止"渠道同时挂在 antigravity / anthropic 两个平台的分组上，用户只访问
//     antigravity，却看到 anthropic 模型"这类跨平台信息泄漏；
//  4. 字段白名单：仅返回用户需要的字段（省略 BillingModelSource / RestrictModels
//     / 内部 ID / Status 等管理字段）。
type AvailableChannelHandler struct {
	channelService *service.ChannelService
	apiKeyService  *service.APIKeyService
	settingService *service.SettingService
}

// NewAvailableChannelHandler 创建用户侧可用渠道 handler。
func NewAvailableChannelHandler(
	channelService *service.ChannelService,
	apiKeyService *service.APIKeyService,
	settingService *service.SettingService,
) *AvailableChannelHandler {
	return &AvailableChannelHandler{
		channelService: channelService,
		apiKeyService:  apiKeyService,
		settingService: settingService,
	}
}

// featureEnabled 返回 available-channels 开关是否启用。默认关闭（opt-in）。
func (h *AvailableChannelHandler) featureEnabled(c *gin.Context) bool {
	if h.settingService == nil {
		return false
	}
	return h.settingService.GetAvailableChannelsRuntime(c.Request.Context()).Enabled
}

// userAvailableGroup 用户可见的分组概要（白名单字段）。
//
// 前端据此区分专属 vs 公开分组（IsExclusive）、订阅 vs 标准分组（SubscriptionType，
// 订阅视觉加深），并用 RateMultiplier 作为默认倍率；用户专属倍率前端走
// /groups/rates，和 API 密钥页面保持一致。
type userAvailableGroup struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Platform         string  `json:"platform"`
	SubscriptionType string  `json:"subscription_type"`
	RateMultiplier   float64 `json:"rate_multiplier"`
	IsExclusive      bool    `json:"is_exclusive"`
}

// userSupportedModelPricing 用户可见的定价字段白名单。
type userSupportedModelPricing struct {
	BillingMode      string                   `json:"billing_mode"`
	InputPrice       *float64                 `json:"input_price"`
	OutputPrice      *float64                 `json:"output_price"`
	CacheWritePrice  *float64                 `json:"cache_write_price"`
	CacheReadPrice   *float64                 `json:"cache_read_price"`
	ImageOutputPrice *float64                 `json:"image_output_price"`
	PerRequestPrice  *float64                 `json:"per_request_price"`
	Intervals        []userPricingIntervalDTO `json:"intervals"`
}

// userPricingIntervalDTO 定价区间白名单（去掉内部 ID、SortOrder 等前端不渲染的字段）。
type userPricingIntervalDTO struct {
	MinTokens       int      `json:"min_tokens"`
	MaxTokens       *int     `json:"max_tokens"`
	TierLabel       string   `json:"tier_label,omitempty"`
	InputPrice      *float64 `json:"input_price"`
	OutputPrice     *float64 `json:"output_price"`
	CacheWritePrice *float64 `json:"cache_write_price"`
	CacheReadPrice  *float64 `json:"cache_read_price"`
	PerRequestPrice *float64 `json:"per_request_price"`
}

// userSupportedModel 用户可见的支持模型条目。
type userSupportedModel struct {
	Name     string                     `json:"name"`
	Platform string                     `json:"platform"`
	Pricing  *userSupportedModelPricing `json:"pricing"`
}

type publicModelPricing struct {
	Name            string   `json:"name"`
	Provider        string   `json:"provider"`
	BillingMode     string   `json:"billing_mode"`
	Currency        string   `json:"currency"`
	InputPrice      *float64 `json:"input_price"`
	OutputPrice     *float64 `json:"output_price"`
	CacheReadPrice  *float64 `json:"cache_read_price"`
	CacheWritePrice *float64 `json:"cache_write_price"`
	PerRequestPrice *float64 `json:"per_request_price"`
	Unit            string   `json:"unit"`
	EndpointTypes   []string `json:"endpoint_types"`
	Tags            []string `json:"tags"`
	Description     string   `json:"description"`
}

// userChannelPlatformSection 单渠道内某个平台的子视图：用户可见的分组 + 该平台
// 支持的模型。按 platform 聚合后让前端可以把渠道名作为 row-group 一次渲染，
// 后面的平台行按 sections 顺序铺开。
type userChannelPlatformSection struct {
	Platform        string               `json:"platform"`
	Groups          []userAvailableGroup `json:"groups"`
	SupportedModels []userSupportedModel `json:"supported_models"`
}

// userAvailableChannel 用户可见的渠道条目（白名单字段）。
//
// 每个渠道聚合为一条记录，内嵌 platforms 子数组：每个 section 对应一个平台，
// 包含该平台的 groups 和 supported_models。
type userAvailableChannel struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Platforms   []userChannelPlatformSection `json:"platforms"`
}

// List 列出当前用户可见的「可用渠道」。
// GET /api/v1/channels/available
func (h *AvailableChannelHandler) List(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Feature 未启用时返回空数组（不暴露渠道信息）。检查放在认证之后，
	// 保持与未开关前的 401 行为一致：未登录先 401，登录后再按开关决定。
	if !h.featureEnabled(c) {
		response.Success(c, []userAvailableChannel{})
		return
	}

	userGroups, err := h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	allowedGroupIDs := make(map[int64]struct{}, len(userGroups))
	for i := range userGroups {
		allowedGroupIDs[userGroups[i].ID] = struct{}{}
	}

	channels, err := h.channelService.ListAvailable(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]userAvailableChannel, 0, len(channels))
	for _, ch := range channels {
		if ch.Status != service.StatusActive {
			continue
		}
		visibleGroups := filterUserVisibleGroups(ch.Groups, allowedGroupIDs)
		if len(visibleGroups) == 0 {
			continue
		}
		sections := buildPlatformSections(ch, visibleGroups)
		if len(sections) == 0 {
			continue
		}
		out = append(out, userAvailableChannel{
			Name:        ch.Name,
			Description: ch.Description,
			Platforms:   sections,
		})
	}

	response.Success(c, out)
}

// ListPublicModels exposes a sanitized, read-only model catalog for the public
// model marketplace. It never returns channel IDs, upstream credentials, base
// URLs, group authorization rules, or any other internal routing data.
func (h *AvailableChannelHandler) ListPublicModels(c *gin.Context) {
	channels, err := h.channelService.ListAvailable(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	byName := make(map[string]publicModelPricing)
	addModel := func(rawName, platform, channelName string, pricing *service.ChannelModelPricing) {
		name := service.CanonicalHuosanyunModelName(rawName)
		if name == "" || strings.Contains(name, "*") || service.IsExcludedHuosanyunModel(name) {
			return
		}
		key := strings.ToLower(name)
		item, exists := byName[key]
		if !exists {
			item = publicModelPricing{
				Name:          name,
				Provider:      publicModelProvider(name, platform, channelName),
				BillingMode:   string(service.BillingModeToken),
				Currency:      "THB",
				Unit:          "1M Tokens",
				EndpointTypes: publicModelEndpointTypes(name),
				Tags:          publicModelTags(name),
				Description:   publicModelDescription(name),
			}
		}
		mergePublicPricing(&item, pricing)
		byName[key] = item
	}

	for _, ch := range channels {
		if ch.Status != service.StatusActive {
			continue
		}
		for _, model := range ch.SupportedModels {
			addModel(model.Name, model.Platform, ch.Name, model.Pricing)
		}
	}

	// Catalog entries are only a non-secret fallback for models that do not exist
	// in channel pricing yet. Do not fill blank channel prices from static data;
	// an operator may intentionally leave a manual price empty.
	for _, pricing := range service.HuosanyunCatalogPricing() {
		if len(pricing.Models) == 0 {
			continue
		}
		name := service.CanonicalHuosanyunModelName(pricing.Models[0])
		if name == "" || service.IsExcludedHuosanyunModel(name) {
			continue
		}
		if _, exists := byName[strings.ToLower(name)]; exists {
			continue
		}
		pricingCopy := pricing.Clone()
		addModel(name, pricing.Platform, "Huosanyun", &pricingCopy)
	}

	out := make([]publicModelPricing, 0, len(byName))
	for _, item := range byName {
		out = append(out, item)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Provider == out[j].Provider {
			return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
		}
		return out[i].Provider < out[j].Provider
	})

	response.Success(c, out)
}
// buildPlatformSections 把一个渠道按 visibleGroups 的平台集合拆成有序的 section 列表：
// 每个 section 对应一个平台，只包含该平台的 groups 和 supported_models。
// 输出按 platform 字母序稳定排序，便于前端等效比较与回归测试。
func buildPlatformSections(
	ch service.AvailableChannel,
	visibleGroups []userAvailableGroup,
) []userChannelPlatformSection {
	groupsByPlatform := make(map[string][]userAvailableGroup, 4)
	for _, g := range visibleGroups {
		if g.Platform == "" {
			continue
		}
		groupsByPlatform[g.Platform] = append(groupsByPlatform[g.Platform], g)
	}
	if len(groupsByPlatform) == 0 {
		return nil
	}

	platforms := make([]string, 0, len(groupsByPlatform))
	for p := range groupsByPlatform {
		platforms = append(platforms, p)
	}
	sort.Strings(platforms)

	sections := make([]userChannelPlatformSection, 0, len(platforms))
	for _, platform := range platforms {
		platformSet := map[string]struct{}{platform: {}}
		sections = append(sections, userChannelPlatformSection{
			Platform:        platform,
			Groups:          groupsByPlatform[platform],
			SupportedModels: toUserSupportedModels(ch.SupportedModels, platformSet),
		})
	}
	return sections
}

// filterUserVisibleGroups 仅保留用户可访问的分组。
func filterUserVisibleGroups(
	groups []service.AvailableGroupRef,
	allowed map[int64]struct{},
) []userAvailableGroup {
	visible := make([]userAvailableGroup, 0, len(groups))
	for _, g := range groups {
		if _, ok := allowed[g.ID]; !ok {
			continue
		}
		visible = append(visible, userAvailableGroup{
			ID:               g.ID,
			Name:             g.Name,
			Platform:         g.Platform,
			SubscriptionType: g.SubscriptionType,
			RateMultiplier:   g.RateMultiplier,
			IsExclusive:      g.IsExclusive,
		})
	}
	return visible
}

// toUserSupportedModels 将 service 层支持模型转换为用户 DTO（字段白名单）。
// 仅保留平台在 allowedPlatforms 中的条目，防止跨平台模型信息泄漏。
// allowedPlatforms 为 nil 时不做平台过滤（保留全部，供测试或明确无过滤场景使用）。
func toUserSupportedModels(
	src []service.SupportedModel,
	allowedPlatforms map[string]struct{},
) []userSupportedModel {
	out := make([]userSupportedModel, 0, len(src))
	for i := range src {
		m := src[i]
		if allowedPlatforms != nil {
			if _, ok := allowedPlatforms[m.Platform]; !ok {
				continue
			}
		}
		out = append(out, userSupportedModel{
			Name:     m.Name,
			Platform: m.Platform,
			Pricing:  toUserPricing(m.Pricing),
		})
	}
	return out
}

// toUserPricing 将 service 层定价转换为用户 DTO；入参为 nil 时返回 nil。
func toUserPricing(p *service.ChannelModelPricing) *userSupportedModelPricing {
	if p == nil {
		return nil
	}
	intervals := make([]userPricingIntervalDTO, 0, len(p.Intervals))
	for _, iv := range p.Intervals {
		intervals = append(intervals, userPricingIntervalDTO{
			MinTokens:       iv.MinTokens,
			MaxTokens:       iv.MaxTokens,
			TierLabel:       iv.TierLabel,
			InputPrice:      iv.InputPrice,
			OutputPrice:     iv.OutputPrice,
			CacheWritePrice: iv.CacheWritePrice,
			CacheReadPrice:  iv.CacheReadPrice,
			PerRequestPrice: iv.PerRequestPrice,
		})
	}
	billingMode := string(p.BillingMode)
	if billingMode == "" {
		billingMode = string(service.BillingModeToken)
	}
	return &userSupportedModelPricing{
		BillingMode:      billingMode,
		InputPrice:       p.InputPrice,
		OutputPrice:      p.OutputPrice,
		CacheWritePrice:  p.CacheWritePrice,
		CacheReadPrice:   p.CacheReadPrice,
		ImageOutputPrice: p.ImageOutputPrice,
		PerRequestPrice:  p.PerRequestPrice,
		Intervals:        intervals,
	}
}

func mergePublicPricing(dst *publicModelPricing, p *service.ChannelModelPricing) {
	if p == nil {
		return
	}
	mode := string(p.BillingMode)
	if mode == "" {
		mode = string(service.BillingModeToken)
	}
	if dst.BillingMode == "" || dst.BillingMode == string(service.BillingModeToken) {
		dst.BillingMode = mode
	}
	if mode == string(service.BillingModePerRequest) || mode == string(service.BillingModeImage) || p.PerRequestPrice != nil {
		dst.Unit = "request"
	} else {
		dst.Unit = "1M Tokens"
	}
	if dst.InputPrice == nil && p.InputPrice != nil {
		dst.InputPrice = p.InputPrice
	}
	if dst.OutputPrice == nil && p.OutputPrice != nil {
		dst.OutputPrice = p.OutputPrice
	}
	if dst.CacheReadPrice == nil && p.CacheReadPrice != nil {
		dst.CacheReadPrice = p.CacheReadPrice
	}
	if dst.CacheWritePrice == nil && p.CacheWritePrice != nil {
		dst.CacheWritePrice = p.CacheWritePrice
	}
	if dst.PerRequestPrice == nil && p.PerRequestPrice != nil {
		dst.PerRequestPrice = p.PerRequestPrice
	}
	for _, interval := range p.Intervals {
		if dst.InputPrice == nil && interval.InputPrice != nil {
			dst.InputPrice = interval.InputPrice
		}
		if dst.OutputPrice == nil && interval.OutputPrice != nil {
			dst.OutputPrice = interval.OutputPrice
		}
		if dst.CacheReadPrice == nil && interval.CacheReadPrice != nil {
			dst.CacheReadPrice = interval.CacheReadPrice
		}
		if dst.CacheWritePrice == nil && interval.CacheWritePrice != nil {
			dst.CacheWritePrice = interval.CacheWritePrice
		}
		if dst.PerRequestPrice == nil && interval.PerRequestPrice != nil {
			dst.PerRequestPrice = interval.PerRequestPrice
		}
	}
}

func publicModelProvider(modelName, platform, channelName string) string {
	name := strings.ToLower(modelName)
	switch {
	case strings.Contains(name, "deepseek"):
		return "DeepSeek"
	case strings.Contains(name, "qwen"):
		return "Qwen"
	case strings.Contains(name, "glm"):
		return "Zhipu/GLM"
	case strings.Contains(name, "kimi") || strings.Contains(name, "moonshot"):
		return "Kimi/Moonshot"
	case strings.Contains(name, "minimax"):
		return "MiniMax"
	case strings.Contains(name, "kling"):
		return "Kling"
	case strings.Contains(name, "doubao"):
		return "Doubao/Seedance"
	case strings.Contains(name, "gpt") || strings.Contains(name, "openai"):
		return "OpenAI"
	}
	if platform != "" {
		return titleASCII(platform)
	}
	if channelName != "" {
		return channelName
	}
	return "Other"
}

func publicModelTags(modelName string) []string {
	name := strings.ToLower(modelName)
	tags := []string{"Token billing"}
	if !strings.Contains(name, "seedance") && !strings.Contains(name, "kling") {
		tags = append(tags, "OpenAI-compatible")
	}
	switch {
	case strings.Contains(name, "deepseek"):
		tags = append(tags, "Reasoning", "Coding")
	case strings.Contains(name, "qwen"):
		tags = append(tags, "Fast response", "General chat")
	case strings.Contains(name, "glm"):
		tags = append(tags, "Text generation", "Tool use")
	case strings.Contains(name, "kimi"):
		tags = append(tags, "Long context", "Office analysis")
	case strings.Contains(name, "kling"):
		tags = append(tags, "Video generation", "Per request")
	case strings.Contains(name, "doubao") || strings.Contains(name, "seedance"):
		tags = append(tags, "Video endpoint", "Content generation")
	case strings.Contains(name, "minimax"):
		tags = append(tags, "Text generation", "Long-form writing")
	default:
		tags = append(tags, "General model")
	}
	return uniqueStrings(tags)
}

func publicModelEndpointTypes(modelName string) []string {
	name := strings.ToLower(modelName)
	switch {
	case strings.Contains(name, "kling"):
		return []string{"video"}
	case strings.Contains(name, "seedance"):
		return []string{"video"}
	default:
		return []string{"openai:/v1/chat/completions"}
	}
}

func publicModelDescription(modelName string) string {
	name := strings.ToLower(modelName)
	switch {
	case strings.Contains(name, "deepseek"):
		return "DeepSeek model for general chat, complex reasoning, code generation, and enterprise text processing."
	case strings.Contains(name, "qwen"):
		return "Qwen model for daily chat, content creation, knowledge Q&A, and low-latency business calls."
	case strings.Contains(name, "glm"):
		return "GLM model for text generation, tool use, office automation, and structured tasks."
	case strings.Contains(name, "kimi"):
		return "Kimi model for long-context understanding, retrieval-augmented workflows, office analysis, and agent tasks."
	case strings.Contains(name, "minimax"):
		return "MiniMax model for general chat, long-form writing, creative content generation, and business integration."
	case strings.Contains(name, "kling"):
		return "Kling video generation model for text-to-video, image-to-video, and creative production."
	case strings.Contains(name, "doubao"):
		return "Doubao model for content generation, multimodal understanding, and high-concurrency business scenarios."
	default:
		return "Model available through this gateway. Actual calls depend on API key permissions, account balance, and backend pricing."
	}
}

func titleASCII(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "Other"
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func uniqueStrings(src []string) []string {
	seen := make(map[string]struct{}, len(src))
	out := make([]string, 0, len(src))
	for _, value := range src {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
