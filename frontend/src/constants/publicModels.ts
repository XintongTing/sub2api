export interface PublicModelInfo {
  id: string
  name: string
  displayName: string
  provider: string
  upstreamModel: string
  type: string
  billing: string
  billingMode: string
  currency?: string
  inputPrice?: number | null
  outputPrice?: number | null
  cacheReadPrice?: number | null
  cacheWritePrice?: number | null
  perRequestPrice?: number | null
  unit: string
  endpointTypes: string[]
  description: string
  descriptionI18n?: Partial<Record<PublicModelLocale, string>>
  tags: string[]
}

export type PublicModelLocale = 'zh-CN' | 'zh-TW' | 'en' | 'th'

const perMillion = (value: number) => value / 1_000_000

export const ALL_PROVIDERS_LABEL = 'All providers'

const descriptions: Record<string, Partial<Record<PublicModelLocale, string>>> = {
  deepseek: {
    'zh-CN': 'DeepSeek 模型，适合通用对话、复杂推理、代码生成和企业文本处理。',
    'zh-TW': 'DeepSeek 模型，適合通用對話、複雜推理、程式碼生成和企業文本處理。',
    en: 'DeepSeek model for general chat, complex reasoning, code generation, and enterprise text processing.',
    th: 'โมเดล DeepSeek สำหรับแชตทั่วไป การให้เหตุผลที่ซับซ้อน การเขียนโค้ด และงานประมวลผลข้อความระดับองค์กร',
  },
  glm: {
    'zh-CN': 'GLM 模型，适合文本生成、工具调用、办公自动化和结构化任务。',
    'zh-TW': 'GLM 模型，適合文本生成、工具調用、辦公自動化和結構化任務。',
    en: 'GLM model for text generation, tool use, office automation, and structured tasks.',
    th: 'โมเดล GLM สำหรับสร้างข้อความ เรียกใช้เครื่องมือ งานอัตโนมัติในสำนักงาน และงานแบบมีโครงสร้าง',
  },
  kimi: {
    'zh-CN': 'Kimi 模型，适合长上下文理解、检索增强、办公分析和 Agent 任务。',
    'zh-TW': 'Kimi 模型，適合長上下文理解、檢索增強、辦公分析和 Agent 任務。',
    en: 'Kimi model for long-context understanding, retrieval-augmented workflows, office analysis, and agent tasks.',
    th: 'โมเดล Kimi สำหรับบริบทยาว เวิร์กโฟลว์แบบเสริมการค้นคืน การวิเคราะห์งานเอกสาร และงาน Agent',
  },
  qwen: {
    'zh-CN': 'Qwen 模型，适合日常对话、内容创作、知识问答和低延迟业务调用。',
    'zh-TW': 'Qwen 模型，適合日常對話、內容創作、知識問答和低延遲業務調用。',
    en: 'Qwen model for daily chat, content creation, knowledge Q&A, and low-latency business calls.',
    th: 'โมเดล Qwen สำหรับแชตประจำวัน การสร้างคอนเทนต์ ถามตอบความรู้ และงานที่ต้องการความหน่วงต่ำ',
  },
  minimax: {
    'zh-CN': 'MiniMax 模型，适合通用对话、长文本写作、创意内容生成和业务集成。',
    'zh-TW': 'MiniMax 模型，適合通用對話、長文本寫作、創意內容生成和業務整合。',
    en: 'MiniMax model for general chat, long-form writing, creative content generation, and business integration.',
    th: 'โมเดล MiniMax สำหรับแชตทั่วไป การเขียนข้อความยาว การสร้างคอนเทนต์เชิงสร้างสรรค์ และการเชื่อมต่อธุรกิจ',
  },
  seedance: {
    'zh-CN': 'Doubao Seedance 模型，面向视频或多模态内容生成场景，具体能力取决于上游端点。',
    'zh-TW': 'Doubao Seedance 模型，面向影片或多模態內容生成場景，具體能力取決於上游端點。',
    en: 'Doubao Seedance model for video or multimodal content generation, depending on the upstream endpoint capability.',
    th: 'โมเดล Doubao Seedance สำหรับสร้างวิดีโอหรือคอนเทนต์หลายรูปแบบ โดยขึ้นอยู่กับความสามารถของปลายทางต้นทาง',
  },
}

function descriptionFor(kind: keyof typeof descriptions): string {
  return descriptions[kind].en || ''
}

export const publicModels: PublicModelInfo[] = [
  {
    id: 'deepseek-v4-flash',
    name: 'deepseek-v4-flash',
    displayName: 'deepseek-v4-flash',
    provider: 'DeepSeek',
    upstreamModel: 'deepseek-v4-flash',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(15),
    outputPrice: perMillion(30),
    cacheReadPrice: perMillion(0.3),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('deepseek'),
    descriptionI18n: descriptions.deepseek,
    tags: ['OpenAI-compatible', 'Token billing', 'Reasoning', 'Fast response'],
  },
  {
    id: 'deepseek-v4-pro',
    name: 'deepseek-v4-pro',
    displayName: 'deepseek-v4-pro',
    provider: 'DeepSeek',
    upstreamModel: 'deepseek-v4-pro',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(180),
    outputPrice: perMillion(360),
    cacheReadPrice: perMillion(36),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('deepseek'),
    descriptionI18n: descriptions.deepseek,
    tags: ['OpenAI-compatible', 'Token billing', 'Reasoning', 'Coding'],
  },
  {
    id: 'doubao-seedance-2-0-fast-idle-260128',
    name: 'doubao-seedance-2-0-fast-idle-260128',
    displayName: 'doubao-seedance-2-0-fast-idle-260128',
    provider: 'Doubao/Seedance',
    upstreamModel: 'doubao-seedance-2-0-fast-idle-260128',
    type: 'Video',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(555),
    outputPrice: perMillion(555),
    unit: '1M Tokens',
    endpointTypes: ['video'],
    description: descriptionFor('seedance'),
    descriptionI18n: descriptions.seedance,
    tags: ['Token billing', 'Video endpoint', 'Content generation'],
  },
  {
    id: 'doubao-seedance-2-0-idle-260128',
    name: 'doubao-seedance-2-0-idle-260128',
    displayName: 'doubao-seedance-2-0-idle-260128',
    provider: 'Doubao/Seedance',
    upstreamModel: 'doubao-seedance-2-0-idle-260128',
    type: 'Video',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(765),
    outputPrice: perMillion(765),
    unit: '1M Tokens',
    endpointTypes: ['video'],
    description: descriptionFor('seedance'),
    descriptionI18n: descriptions.seedance,
    tags: ['Token billing', 'Video endpoint', 'Content generation'],
  },
  {
    id: 'glm-4.7',
    name: 'glm-4.7',
    displayName: 'glm-4.7',
    provider: 'Zhipu/GLM',
    upstreamModel: 'glm-4.7',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(4.05),
    outputPrice: perMillion(16.5),
    cacheReadPrice: perMillion(8.22),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('glm'),
    descriptionI18n: descriptions.glm,
    tags: ['OpenAI-compatible', 'Token billing', 'Text generation', 'Tool use'],
  },
  {
    id: 'glm-5',
    name: 'glm-5',
    displayName: 'glm-5',
    provider: 'Zhipu/GLM',
    upstreamModel: 'glm-5',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(11.25),
    outputPrice: perMillion(36),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('glm'),
    descriptionI18n: descriptions.glm,
    tags: ['OpenAI-compatible', 'Token billing', 'Text generation', 'Tool use'],
  },
  {
    id: 'glm-5.1',
    name: 'glm-5.1',
    displayName: 'glm-5.1',
    provider: 'Zhipu/GLM',
    upstreamModel: 'glm-5.1',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(90),
    outputPrice: perMillion(360),
    cacheReadPrice: perMillion(18),
    cacheWritePrice: perMillion(112.5),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('glm'),
    descriptionI18n: descriptions.glm,
    tags: ['OpenAI-compatible', 'Token billing', 'Coding', 'Tool use'],
  },
  {
    id: 'kimi-k2.5',
    name: 'kimi-k2.5',
    displayName: 'kimi-k2.5',
    provider: 'Kimi/Moonshot',
    upstreamModel: 'kimi-k2.5',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(8.25),
    outputPrice: perMillion(41.4),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('kimi'),
    descriptionI18n: descriptions.kimi,
    tags: ['OpenAI-compatible', 'Token billing', 'Long context', 'Office analysis'],
  },
  {
    id: 'kimi-k2.6',
    name: 'kimi-k2.6',
    displayName: 'kimi-k2.6',
    provider: 'Kimi/Moonshot',
    upstreamModel: 'kimi-k2.6',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(97.5),
    outputPrice: perMillion(405),
    cacheReadPrice: perMillion(15),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('kimi'),
    descriptionI18n: descriptions.kimi,
    tags: ['OpenAI-compatible', 'Token billing', 'Long context', 'Office analysis'],
  },
  {
    id: 'MiniMax-M2.5',
    name: 'MiniMax-M2.5',
    displayName: 'MiniMax-M2.5',
    provider: 'MiniMax',
    upstreamModel: 'MiniMax-M2.5',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: null,
    outputPrice: null,
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('minimax'),
    descriptionI18n: descriptions.minimax,
    tags: ['OpenAI-compatible', 'Token billing', 'Text generation', 'Long-form writing'],
  },
  {
    id: 'qwen-plus',
    name: 'qwen-plus',
    displayName: 'qwen-plus',
    provider: 'Qwen',
    upstreamModel: 'qwen-plus',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: null,
    outputPrice: null,
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('qwen'),
    descriptionI18n: descriptions.qwen,
    tags: ['OpenAI-compatible', 'Token billing', 'General chat', 'Fast response'],
  },
  {
    id: 'qwen3.6-plus',
    name: 'qwen3.6-plus',
    displayName: 'qwen3.6-plus',
    provider: 'Qwen',
    upstreamModel: 'qwen3.6-plus',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: null,
    outputPrice: null,
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: descriptionFor('qwen'),
    descriptionI18n: descriptions.qwen,
    tags: ['OpenAI-compatible', 'Token billing', 'General chat', 'Reasoning'],
  },
]


export function normalizePublicModelLocale(locale: string): PublicModelLocale {
  const value = locale.toLowerCase()
  if (value.startsWith('zh-tw')) return 'zh-TW'
  if (value.startsWith('en')) return 'en'
  if (value.startsWith('th')) return 'th'
  return 'zh-CN'
}

export function canonicalPublicModelName(name: string): string {
  const trimmed = String(name || '').trim()
  if (!trimmed) return ''
  return trimmed
}

export function isPublicModelExcluded(name: string): boolean {
  const normalized = String(name || '').trim().toLowerCase()
  return (
    normalized === 'deepseek-v3.2' ||
    normalized === 'deepseek-v3' ||
    normalized === 'deepseek-pro' ||
    normalized === 'qwen3-turbo' ||
    normalized === 'glm4-air' ||
    normalized === 'glm4-plus' ||
    normalized.startsWith('kling')
  )
}

export function dedupePublicModels(models: PublicModelInfo[]): PublicModelInfo[] {
  const seen = new Set<string>()
  const output: PublicModelInfo[] = []

  for (const model of models) {
    const rawValues = [model.id, model.name, model.displayName, model.upstreamModel].filter(Boolean) as string[]
    if (rawValues.some(isPublicModelExcluded)) continue

    const canonicalName = canonicalPublicModelName(model.name || model.id || model.displayName)
    if (!canonicalName || isPublicModelExcluded(canonicalName)) continue

    const key = canonicalName.toLowerCase()
    if (seen.has(key)) continue
    seen.add(key)

    output.push({
      ...model,
      id: canonicalName,
      name: canonicalName,
      displayName: canonicalName,
      upstreamModel: canonicalName,
      currency: model.currency || 'THB',
      tags: Array.from(new Set(model.tags || [])),
    })
  }

  return output.sort((a, b) => a.name.localeCompare(b.name))
}

export function findPublicModelByName(name: string): PublicModelInfo | undefined {
  const canonicalName = canonicalPublicModelName(name).toLowerCase()
  return publicModels.find(model => model.name.toLowerCase() === canonicalName)
}

export function localizePublicModelDescription(model: PublicModelInfo, locale: string): string {
  const key = normalizePublicModelLocale(locale)
  return model.descriptionI18n?.[key] || model.descriptionI18n?.en || model.description || ''
}

const tagLabels: Record<string, Partial<Record<PublicModelLocale, string>>> = {
  'OpenAI-compatible': {
    'zh-CN': 'OpenAI 兼容',
    'zh-TW': 'OpenAI 相容',
    en: 'OpenAI-compatible',
    th: 'รองรับ OpenAI',
  },
  'Token billing': {
    'zh-CN': '按量计费',
    'zh-TW': '按量計費',
    en: 'Token billing',
    th: 'คิดตามโทเคน',
  },
  'Reasoning': {
    'zh-CN': '推理',
    'zh-TW': '推理',
    en: 'Reasoning',
    th: 'เหตุผล',
  },
  'Fast response': {
    'zh-CN': '快速响应',
    'zh-TW': '快速回應',
    en: 'Fast response',
    th: 'ตอบสนองเร็ว',
  },
  'Coding': {
    'zh-CN': '代码',
    'zh-TW': '程式碼',
    en: 'Coding',
    th: 'เขียนโค้ด',
  },
  'Text generation': {
    'zh-CN': '文本生成',
    'zh-TW': '文本生成',
    en: 'Text generation',
    th: 'สร้างข้อความ',
  },
  'Tool use': {
    'zh-CN': '工具调用',
    'zh-TW': '工具調用',
    en: 'Tool use',
    th: 'ใช้เครื่องมือ',
  },
  'Video endpoint': {
    'zh-CN': '视频端点',
    'zh-TW': '影片端點',
    en: 'Video endpoint',
    th: 'ปลายทางวิดีโอ',
  },
  'Content generation': {
    'zh-CN': '内容生成',
    'zh-TW': '內容生成',
    en: 'Content generation',
    th: 'สร้างคอนเทนต์',
  },
  'Long context': {
    'zh-CN': '长上下文',
    'zh-TW': '長上下文',
    en: 'Long context',
    th: 'บริบทยาว',
  },
  'Office analysis': {
    'zh-CN': '办公分析',
    'zh-TW': '辦公分析',
    en: 'Office analysis',
    th: 'วิเคราะห์งานเอกสาร',
  },
  'Long-form writing': {
    'zh-CN': '长文写作',
    'zh-TW': '長文寫作',
    en: 'Long-form writing',
    th: 'เขียนข้อความยาว',
  },
  'General chat': {
    'zh-CN': '通用对话',
    'zh-TW': '通用對話',
    en: 'General chat',
    th: 'แชตทั่วไป',
  },
}

export function localizePublicModelTag(tag: string, locale: string): string {
  const key = normalizePublicModelLocale(locale)
  return tagLabels[tag]?.[key] || tag
}

export function localizePublicModelTags(tags: string[], locale: string): string[] {
  return tags.map(tag => localizePublicModelTag(tag, locale))
}
