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

export const publicModels: PublicModelInfo[] = [
  {
    id: 'qwen-plus',
    name: 'qwen-plus',
    displayName: 'Qwen3-Turbo',
    provider: 'Qwen',
    upstreamModel: 'qwen-plus',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(8),
    outputPrice: perMillion(24),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: 'Fast general chat model for customer support, content generation, knowledge Q&A, and low-latency business calls.',
    descriptionI18n: {
      'zh-CN': '快速通用对话模型，适合客服、内容生成、知识问答和低延迟业务调用。',
      'zh-TW': '快速通用對話模型，適合客服、內容生成、知識問答和低延遲業務調用。',
      en: 'Fast general chat model for customer support, content generation, knowledge Q&A, and low-latency business calls.',
      th: 'โมเดลสนทนาทั่วไปที่ตอบสนองรวดเร็ว เหมาะกับงานบริการลูกค้า การสร้างเนื้อหา ถามตอบความรู้ และงานธุรกิจที่ต้องการความหน่วงต่ำ',
    },
    tags: ['Fast response', 'General chat', 'Cost effective'],
  },
  {
    id: 'deepseek-v3.2',
    name: 'deepseek-v3.2',
    displayName: 'DeepSeek-V3',
    provider: 'DeepSeek',
    upstreamModel: 'deepseek-v3.2',
    type: 'Chat',
    billing: 'Token billing',
    billingMode: 'token',
    currency: 'THB',
    inputPrice: perMillion(30),
    outputPrice: perMillion(45),
    cacheReadPrice: perMillion(6),
    cacheWritePrice: perMillion(37.5),
    unit: '1M Tokens',
    endpointTypes: ['openai:/v1/chat/completions'],
    description: 'General-purpose DeepSeek model for enterprise text processing, code assistance, reasoning, and complex Q&A.',
    descriptionI18n: {
      'zh-CN': '通用 DeepSeek 模型，适合企业文本处理、代码辅助、推理和复杂问答。',
      'zh-TW': '通用 DeepSeek 模型，適合企業文本處理、程式碼輔助、推理和複雜問答。',
      en: 'General-purpose DeepSeek model for enterprise text processing, code assistance, reasoning, and complex Q&A.',
      th: 'โมเดล DeepSeek สำหรับงานทั่วไป เหมาะกับการประมวลผลข้อความองค์กร ผู้ช่วยเขียนโค้ด การให้เหตุผล และคำถามซับซ้อน',
    },
    tags: ['Reasoning', 'Coding', 'General chat'],
  },
  {
    id: 'glm-4.7',
    name: 'glm-4.7',
    displayName: 'GLM4-Air',
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
    description: 'Lightweight GLM model for text generation, customer chat, structured office tasks, and routine automation.',
    descriptionI18n: {
      'zh-CN': '轻量级 GLM 模型，适合文本生成、客户对话、结构化办公任务和日常自动化。',
      'zh-TW': '輕量級 GLM 模型，適合文字生成、客戶對話、結構化辦公任務和日常自動化。',
      en: 'Lightweight GLM model for text generation, customer chat, structured office tasks, and routine automation.',
      th: 'โมเดล GLM น้ำหนักเบา เหมาะกับการสร้างข้อความ แชทลูกค้า งานเอกสารแบบมีโครงสร้าง และระบบอัตโนมัติทั่วไป',
    },
    tags: ['Text generation', 'Office automation', 'Lightweight'],
  },
  {
    id: 'glm-5',
    name: 'glm-5',
    displayName: 'GLM4-Plus',
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
    description: 'Enhanced GLM model for stronger text generation, reasoning, tool use, and business workflow scenarios.',
    descriptionI18n: {
      'zh-CN': '增强型 GLM 模型，适合更强文本生成、推理、工具调用和业务流程场景。',
      'zh-TW': '增強型 GLM 模型，適合更強文字生成、推理、工具調用和業務流程場景。',
      en: 'Enhanced GLM model for stronger text generation, reasoning, tool use, and business workflow scenarios.',
      th: 'โมเดล GLM รุ่นเสริม เหมาะกับการสร้างข้อความที่ซับซ้อนขึ้น การให้เหตุผล การใช้เครื่องมือ และเวิร์กโฟลว์ธุรกิจ',
    },
    tags: ['Advanced reasoning', 'Tool use', 'Business workflows'],
  },
  {
    id: 'deepseek-v4-pro',
    name: 'deepseek-v4-pro',
    displayName: 'DeepSeek-Pro',
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
    description: 'Premium DeepSeek model for advanced reasoning, professional coding, research analysis, and long-context understanding.',
    descriptionI18n: {
      'zh-CN': '高阶 DeepSeek 模型，适合复杂推理、专业代码、研究分析和长上下文理解。',
      'zh-TW': '高階 DeepSeek 模型，適合複雜推理、專業程式碼、研究分析和長上下文理解。',
      en: 'Premium DeepSeek model for advanced reasoning, professional coding, research analysis, and long-context understanding.',
      th: 'โมเดล DeepSeek ระดับสูง เหมาะกับการให้เหตุผลขั้นสูง งานโค้ดมืออาชีพ การวิเคราะห์เชิงวิจัย และบริบทยาว',
    },
    tags: ['Pro capability', 'Complex reasoning', 'Long context'],
  },
]

const tagLabels: Record<string, Record<PublicModelLocale, string>> = {
  'Fast response': {
    'zh-CN': '快速响应',
    'zh-TW': '快速回應',
    en: 'Fast response',
    th: 'ตอบสนองรวดเร็ว',
  },
  'General chat': {
    'zh-CN': '通用对话',
    'zh-TW': '通用對話',
    en: 'General chat',
    th: 'สนทนาทั่วไป',
  },
  'Cost effective': {
    'zh-CN': '高性价比',
    'zh-TW': '高性價比',
    en: 'Cost effective',
    th: 'คุ้มค่า',
  },
  Reasoning: {
    'zh-CN': '推理',
    'zh-TW': '推理',
    en: 'Reasoning',
    th: 'การให้เหตุผล',
  },
  Coding: {
    'zh-CN': '代码',
    'zh-TW': '程式碼',
    en: 'Coding',
    th: 'เขียนโค้ด',
  },
  'Text generation': {
    'zh-CN': '文本生成',
    'zh-TW': '文字生成',
    en: 'Text generation',
    th: 'สร้างข้อความ',
  },
  'Office automation': {
    'zh-CN': '办公自动化',
    'zh-TW': '辦公自動化',
    en: 'Office automation',
    th: 'งานสำนักงานอัตโนมัติ',
  },
  Lightweight: {
    'zh-CN': '轻量',
    'zh-TW': '輕量',
    en: 'Lightweight',
    th: 'เบาและเร็ว',
  },
  'Advanced reasoning': {
    'zh-CN': '高级推理',
    'zh-TW': '高階推理',
    en: 'Advanced reasoning',
    th: 'เหตุผลขั้นสูง',
  },
  'Tool use': {
    'zh-CN': '工具调用',
    'zh-TW': '工具調用',
    en: 'Tool use',
    th: 'ใช้เครื่องมือ',
  },
  'Business workflows': {
    'zh-CN': '业务流程',
    'zh-TW': '業務流程',
    en: 'Business workflows',
    th: 'เวิร์กโฟลว์ธุรกิจ',
  },
  'Pro capability': {
    'zh-CN': '专业能力',
    'zh-TW': '專業能力',
    en: 'Pro capability',
    th: 'ความสามารถระดับโปร',
  },
  'Complex reasoning': {
    'zh-CN': '复杂推理',
    'zh-TW': '複雜推理',
    en: 'Complex reasoning',
    th: 'เหตุผลซับซ้อน',
  },
  'Long context': {
    'zh-CN': '长上下文',
    'zh-TW': '長上下文',
    en: 'Long context',
    th: 'บริบทยาว',
  },
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
    th: 'คิดเงินตาม Token',
  },
}

export function normalizePublicModelLocale(locale: string): PublicModelLocale {
  if (locale.startsWith('zh-TW')) return 'zh-TW'
  if (locale.startsWith('en')) return 'en'
  if (locale.startsWith('th')) return 'th'
  return 'zh-CN'
}

export function localizePublicModelDescription(model: PublicModelInfo, locale: string): string {
  const key = normalizePublicModelLocale(locale)
  const known = findPublicModelByName(model.name) || findPublicModelByName(model.displayName)
  return model.descriptionI18n?.[key] || known?.descriptionI18n?.[key] || model.description
}

export function localizePublicModelTag(tag: string, locale: string): string {
  const key = normalizePublicModelLocale(locale)
  return tagLabels[tag]?.[key] || tag
}

export function localizePublicModelTags(tags: string[], locale: string): string[] {
  return tags.map(tag => localizePublicModelTag(tag, locale))
}

export function findPublicModelByName(name: string): PublicModelInfo | undefined {
  const normalized = name.trim().toLowerCase()
  return publicModels.find((model) =>
    [model.id, model.name, model.displayName, model.upstreamModel]
      .filter(Boolean)
      .some((value) => value.toLowerCase() === normalized)
  )
}

export function makeProviderOptions(models: PublicModelInfo[]): string[] {
  return [
    ALL_PROVIDERS_LABEL,
    ...Array.from(new Set(models.map((model) => model.provider).filter(Boolean))),
  ]
}
