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
  tags: string[]
}

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
    tags: ['Pro capability', 'Complex reasoning', 'Long context'],
  },
]

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
