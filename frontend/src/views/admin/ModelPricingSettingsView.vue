<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="rounded-lg border border-emerald-100 bg-gradient-to-r from-primary-700 via-primary-600 to-cyan-500 p-6 text-white shadow-sm">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <p class="text-sm font-semibold text-primary-50">财务管理</p>
            <h1 class="mt-2 text-2xl font-bold sm:text-3xl">模型价格设置</h1>
            <p class="mt-2 max-w-3xl text-sm leading-6 text-primary-50">
              直接维护可售模型的 THB 售价。保存后会写入渠道 model_pricing，前台模型广场和新 API 调用扣费都会读取同一套生效价格。
            </p>
          </div>
          <div class="flex flex-wrap gap-3">
            <button class="rounded-md bg-white/15 px-4 py-2 text-sm font-semibold text-white hover:bg-white/25" :disabled="loading" @click="load">
              {{ loading ? '刷新中...' : '刷新' }}
            </button>
            <router-link to="/admin/channels/pricing" class="rounded-md bg-white px-4 py-2 text-sm font-semibold text-primary-700 hover:bg-primary-50">
              高级渠道配置
            </router-link>
          </div>
        </div>
      </div>

      <div class="grid gap-4 md:grid-cols-4">
        <div class="card p-4">
          <p class="text-sm text-slate-500 dark:text-dark-300">可售模型</p>
          <p class="mt-2 text-2xl font-bold text-slate-950 dark:text-white">{{ rows.length }}</p>
        </div>
        <div class="card p-4">
          <p class="text-sm text-slate-500 dark:text-dark-300">已启用渠道</p>
          <p class="mt-2 text-2xl font-bold text-slate-950 dark:text-white">{{ activeChannels.length }}</p>
        </div>
        <div class="card p-4">
          <p class="text-sm text-slate-500 dark:text-dark-300">未保存修改</p>
          <p class="mt-2 text-2xl font-bold text-primary-700 dark:text-primary-300">{{ dirtyRows.length }}</p>
        </div>
        <div class="card p-4">
          <p class="text-sm text-slate-500 dark:text-dark-300">价格单位</p>
          <p class="mt-2 text-lg font-bold text-slate-950 dark:text-white">THB / 泰铢</p>
        </div>
      </div>

      <div v-if="error" class="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-200">
        {{ error }}
      </div>

      <div class="card p-4">
        <div class="grid gap-3 lg:grid-cols-[1fr_180px_180px_auto]">
          <input
            v-model.trim="search"
            type="search"
            placeholder="搜索模型名称 / 供应商 / 渠道"
            class="rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 outline-none focus:border-primary-400 dark:border-dark-700 dark:bg-dark-900 dark:text-white"
          />
          <select v-model="providerFilter" class="rounded-md border border-slate-200 bg-white px-3 py-2 text-sm dark:border-dark-700 dark:bg-dark-900 dark:text-white">
            <option value="">全部供应商</option>
            <option v-for="provider in providers" :key="provider" :value="provider">{{ provider }}</option>
          </select>
          <select v-model="billingFilter" class="rounded-md border border-slate-200 bg-white px-3 py-2 text-sm dark:border-dark-700 dark:bg-dark-900 dark:text-white">
            <option value="">全部计费类型</option>
            <option value="token">按量计费</option>
            <option value="per_request">按次计费</option>
            <option value="image">图片/按次计费</option>
          </select>
          <button
            class="rounded-md bg-primary-600 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="dirtyRows.length === 0 || bulkSaving"
            @click="saveDirtyRows"
          >
            {{ bulkSaving ? '保存中...' : `保存全部 (${dirtyRows.length})` }}
          </button>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-12">
        <LoadingSpinner />
      </div>

      <div v-else class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm dark:border-dark-800 dark:bg-dark-900">
        <div class="overflow-x-auto">
          <table class="min-w-[1320px] w-full divide-y divide-slate-200 text-sm dark:divide-dark-700">
            <thead class="bg-slate-50 text-left text-xs font-semibold uppercase text-slate-500 dark:bg-dark-800 dark:text-dark-300">
              <tr>
                <th class="px-4 py-3">模型</th>
                <th class="px-4 py-3">供应商</th>
                <th class="px-4 py-3">端点</th>
                <th class="px-4 py-3">计费类型</th>
                <th class="px-4 py-3">输入价<br />฿ / 1M Tokens</th>
                <th class="px-4 py-3">补全价<br />฿ / 1M Tokens</th>
                <th class="px-4 py-3">缓存读取<br />฿ / 1M Tokens</th>
                <th class="px-4 py-3">缓存创建<br />฿ / 1M Tokens</th>
                <th class="px-4 py-3">按次价<br />฿ / 次</th>
                <th class="px-4 py-3">前台/调用</th>
                <th class="px-4 py-3">更新时间</th>
                <th class="px-4 py-3 text-right">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 dark:divide-dark-800">
              <tr v-for="row in filteredRows" :key="row.id" class="hover:bg-primary-50/40 dark:hover:bg-primary-500/5">
                <td class="px-4 py-3 align-top">
                  <div class="font-mono font-semibold text-slate-950 dark:text-white">{{ row.model }}</div>
                  <div class="mt-1 text-xs text-slate-500 dark:text-dark-300">{{ row.channelName || '待写入主渠道' }}</div>
                </td>
                <td class="px-4 py-3 align-top text-slate-700 dark:text-dark-200">{{ row.provider }}</td>
                <td class="px-4 py-3 align-top">
                  <span class="rounded-full bg-primary-50 px-2 py-1 text-xs font-semibold text-primary-700 dark:bg-primary-500/10 dark:text-primary-300">
                    {{ row.endpointType }}
                  </span>
                </td>
                <td class="px-4 py-3 align-top">
                  <select v-model="row.billingMode" class="w-32 rounded-md border border-slate-200 bg-white px-2 py-1.5 text-sm dark:border-dark-700 dark:bg-dark-950 dark:text-white" @change="markDirty(row)">
                    <option value="token">按量计费</option>
                    <option value="per_request">按次计费</option>
                    <option value="image">图片/按次</option>
                  </select>
                </td>
                <td class="px-4 py-3 align-top"><PriceInput v-model="row.inputPriceText" @update:model-value="markDirty(row)" /></td>
                <td class="px-4 py-3 align-top"><PriceInput v-model="row.outputPriceText" @update:model-value="markDirty(row)" /></td>
                <td class="px-4 py-3 align-top"><PriceInput v-model="row.cacheReadPriceText" @update:model-value="markDirty(row)" /></td>
                <td class="px-4 py-3 align-top"><PriceInput v-model="row.cacheWritePriceText" @update:model-value="markDirty(row)" /></td>
                <td class="px-4 py-3 align-top"><PriceInput v-model="row.perRequestPriceText" @update:model-value="markDirty(row)" /></td>
                <td class="px-4 py-3 align-top">
                  <div class="flex flex-col gap-2 text-xs text-slate-600 dark:text-dark-200">
                    <label class="inline-flex items-center gap-2">
                      <input v-model="row.publicVisible" type="checkbox" class="h-4 w-4 rounded border-slate-300 text-primary-600 focus:ring-primary-500" @change="markDirty(row)" />
                      <span>前台显示</span>
                    </label>
                    <label class="inline-flex items-center gap-2">
                      <input v-model="row.apiEnabled" type="checkbox" class="h-4 w-4 rounded border-slate-300 text-primary-600 focus:ring-primary-500" @change="markDirty(row)" />
                      <span>允许调用</span>
                    </label>
                  </div>
                </td>
                <td class="px-4 py-3 align-top text-xs text-slate-500 dark:text-dark-300">{{ formatDate(row.updatedAt) }}</td>
                <td class="px-4 py-3 align-top text-right">
                  <div class="flex flex-col items-end gap-2">
                    <button
                      class="rounded-md bg-primary-600 px-3 py-1.5 text-xs font-semibold text-white hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50"
                      :disabled="row.saving || !row.dirty"
                      @click="saveRow(row)"
                    >
                      {{ row.saving ? '保存中...' : '保存' }}
                    </button>
                    <span v-if="row.error" class="max-w-48 text-xs text-red-600">{{ row.error }}</span>
                    <span v-else-if="row.saved" class="text-xs text-emerald-600">已保存</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-if="filteredRows.length === 0" class="p-10 text-center text-sm text-slate-500 dark:text-dark-300">
          暂无可编辑模型。请先在高级渠道配置里同步或添加可售模型。
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import channelsAPI, { type Channel, type ChannelModelPricing } from '@/api/admin/channels'
import { CHANNEL_STATUS_ACTIVE, BILLING_MODE_TOKEN, BILLING_MODE_PER_REQUEST, BILLING_MODE_IMAGE, type BillingMode } from '@/constants/channel'
import { canonicalPublicModelName, dedupePublicModels, isPublicModelExcluded, publicModels, type PublicModelInfo } from '@/constants/publicModels'

const PER_MILLION = 1_000_000

const PriceInput = defineComponent({
  props: {
    modelValue: { type: String, default: '' },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    return () => h('input', {
      value: props.modelValue,
      inputmode: 'decimal',
      placeholder: 'Not set',
      class: 'w-28 rounded-md border border-slate-200 bg-white px-2 py-1.5 text-sm text-slate-900 outline-none focus:border-primary-400 dark:border-dark-700 dark:bg-dark-950 dark:text-white',
      onInput: (event: Event) => emit('update:modelValue', (event.target as HTMLInputElement).value),
    })
  },
})

interface PricingRow {
  id: string
  model: string
  provider: string
  endpointType: string
  billingMode: BillingMode
  channelId: number | null
  channelName: string
  channelStatus: string
  entryIndex: number | null
  entryId?: number
  platform: string
  inputPriceText: string
  outputPriceText: string
  cacheReadPriceText: string
  cacheWritePriceText: string
  perRequestPriceText: string
  publicVisible: boolean
  apiEnabled: boolean
  updatedAt?: string
  source: 'channel' | 'fallback'
  dirty: boolean
  saving: boolean
  saved: boolean
  error: string
}

const loading = ref(false)
const bulkSaving = ref(false)
const error = ref('')
const channels = ref<Channel[]>([])
const rows = ref<PricingRow[]>([])
const search = ref('')
const providerFilter = ref('')
const billingFilter = ref('')

const knownModels = computed(() => {
  const map = new Map<string, PublicModelInfo>()
  for (const model of dedupePublicModels(publicModels)) {
    map.set(canonicalPublicModelName(model.name).toLowerCase(), model)
  }
  return map
})

const activeChannels = computed(() => channels.value.filter(channel => channel.status === CHANNEL_STATUS_ACTIVE))
const primaryChannel = computed(() => activeChannels.value[0] || channels.value[0] || null)
const dirtyRows = computed(() => rows.value.filter(row => row.dirty))
const providers = computed(() => Array.from(new Set(rows.value.map(row => row.provider))).sort())

const filteredRows = computed(() => {
  const q = search.value.trim().toLowerCase()
  return rows.value.filter(row => {
    if (providerFilter.value && row.provider !== providerFilter.value) return false
    if (billingFilter.value && row.billingMode !== billingFilter.value) return false
    if (!q) return true
    return [row.model, row.provider, row.channelName, row.endpointType].some(value => value.toLowerCase().includes(q))
  })
})

function formatTokenPrice(value: number | null | undefined): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) return ''
  return trimNumber(value * PER_MILLION)
}

function formatDirectPrice(value: number | null | undefined): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) return ''
  return trimNumber(value)
}

function trimNumber(value: number): string {
  return value.toFixed(6).replace(/\.0+$/, '').replace(/(\.\d*?)0+$/, '$1')
}

function parseTokenPrice(value: string): number | null {
  const parsed = parsePrice(value)
  return parsed == null ? null : parsed / PER_MILLION
}

function parseDirectPrice(value: string): number | null {
  return parsePrice(value)
}

function parsePrice(value: string): number | null {
  const trimmed = value.trim()
  if (!trimmed) return null
  const parsed = Number(trimmed)
  if (!Number.isFinite(parsed) || parsed < 0) {
    throw new Error('Price must be a number greater than or equal to 0')
  }
  return parsed
}

function clonePricing(entry: ChannelModelPricing): ChannelModelPricing {
  return {
    ...entry,
    models: [...(entry.models || [])],
    intervals: (entry.intervals || []).map(interval => ({ ...interval })),
  }
}

function endpointFor(model: string, known?: PublicModelInfo): string {
  if (known?.endpointTypes?.some(endpoint => endpoint === 'video')) return 'video'
  if (model.toLowerCase().includes('seedance')) return 'video'
  return '/v1/chat/completions'
}

function providerFor(model: string, channel: Channel | null, known?: PublicModelInfo): string {
  if (known?.provider) return known.provider
  const lower = model.toLowerCase()
  if (lower.includes('deepseek')) return 'DeepSeek'
  if (lower.includes('glm')) return 'Zhipu/GLM'
  if (lower.includes('qwen')) return 'Qwen'
  if (lower.includes('kimi')) return 'Kimi'
  if (lower.includes('minimax')) return 'MiniMax'
  if (lower.includes('doubao') || lower.includes('seedance')) return 'Doubao'
  return channel?.name || 'Custom'
}

function rowFromPricing(channel: Channel, entry: ChannelModelPricing, entryIndex: number, modelName: string): PricingRow | null {
  const model = canonicalPublicModelName(modelName)
  if (!model || model.includes('*') || isPublicModelExcluded(model)) return null
  const known = knownModels.value.get(model.toLowerCase())
  return {
    id: `${channel.id}:${entry.id || entryIndex}:${model.toLowerCase()}`,
    model,
    provider: providerFor(model, channel, known),
    endpointType: endpointFor(model, known),
    billingMode: (entry.billing_mode || BILLING_MODE_TOKEN) as BillingMode,
    channelId: channel.id,
    channelName: channel.name,
    channelStatus: channel.status,
    entryIndex,
    entryId: entry.id,
    platform: entry.platform || 'openai',
    inputPriceText: formatTokenPrice(entry.input_price),
    outputPriceText: formatTokenPrice(entry.output_price),
    cacheReadPriceText: formatTokenPrice(entry.cache_read_price),
    cacheWritePriceText: formatTokenPrice(entry.cache_write_price),
    perRequestPriceText: formatDirectPrice(entry.per_request_price),
    publicVisible: entry.public_visible ?? (channel.status === CHANNEL_STATUS_ACTIVE),
    apiEnabled: entry.api_enabled ?? (channel.status === CHANNEL_STATUS_ACTIVE),
    updatedAt: entry.updated_at || channel.updated_at,
    source: 'channel',
    dirty: false,
    saving: false,
    saved: false,
    error: '',
  }
}

function fallbackRow(model: PublicModelInfo): PricingRow {
  return {
    id: `fallback:${model.name.toLowerCase()}`,
    model: canonicalPublicModelName(model.name),
    provider: model.provider,
    endpointType: endpointFor(model.name, model),
    billingMode: (model.billingMode || BILLING_MODE_TOKEN) as BillingMode,
    channelId: primaryChannel.value?.id || null,
    channelName: primaryChannel.value?.name || '',
    channelStatus: primaryChannel.value?.status || '',
    entryIndex: null,
    platform: 'openai',
    inputPriceText: formatTokenPrice(model.inputPrice),
    outputPriceText: formatTokenPrice(model.outputPrice),
    cacheReadPriceText: formatTokenPrice(model.cacheReadPrice),
    cacheWritePriceText: formatTokenPrice(model.cacheWritePrice),
    perRequestPriceText: formatDirectPrice(model.perRequestPrice),
    publicVisible: Boolean(primaryChannel.value && primaryChannel.value.status === CHANNEL_STATUS_ACTIVE),
    apiEnabled: Boolean(primaryChannel.value && primaryChannel.value.status === CHANNEL_STATUS_ACTIVE),
    updatedAt: primaryChannel.value?.updated_at,
    source: 'fallback',
    dirty: false,
    saving: false,
    saved: false,
    error: '',
  }
}

function rebuildRows(): void {
  const byModel = new Map<string, PricingRow>()
  const sortedChannels = [...channels.value].sort((a, b) => {
    if (a.status === b.status) return a.id - b.id
    return a.status === CHANNEL_STATUS_ACTIVE ? -1 : 1
  })

  for (const channel of sortedChannels) {
    (channel.model_pricing || []).forEach((entry, entryIndex) => {
      for (const raw of entry.models || []) {
        const row = rowFromPricing(channel, entry, entryIndex, raw)
        if (!row) continue
        const key = row.model.toLowerCase()
        if (!byModel.has(key)) {
          byModel.set(key, row)
        }
      }
    })
  }

  for (const model of dedupePublicModels(publicModels)) {
    const canonical = canonicalPublicModelName(model.name)
    if (!canonical || isPublicModelExcluded(canonical)) continue
    const key = canonical.toLowerCase()
    if (!byModel.has(key)) {
      byModel.set(key, fallbackRow(model))
    }
  }

  rows.value = Array.from(byModel.values()).sort((a, b) => a.model.localeCompare(b.model))
}

async function load(): Promise<void> {
  loading.value = true
  error.value = ''
  try {
    const loaded: Channel[] = []
    let page = 1
    let total = 0
    do {
      const res = await channelsAPI.list(page, 100, undefined)
      loaded.push(...res.items)
      total = res.total
      page += 1
    } while (loaded.length < total)
    channels.value = loaded
    rebuildRows()
  } catch (err: any) {
    error.value = err?.message || '加载模型价格失败'
  } finally {
    loading.value = false
  }
}

function markDirty(row: PricingRow): void {
  row.dirty = true
  row.saved = false
  row.error = ''
}

function validateRow(row: PricingRow): void {
  parseTokenPrice(row.inputPriceText)
  parseTokenPrice(row.outputPriceText)
  parseTokenPrice(row.cacheReadPriceText)
  parseTokenPrice(row.cacheWritePriceText)
  const perRequest = parseDirectPrice(row.perRequestPriceText)
  if ((row.billingMode === BILLING_MODE_PER_REQUEST || row.billingMode === BILLING_MODE_IMAGE) && perRequest == null) {
    throw new Error('按次计费模型必须填写按次价格')
  }
}

function applyRowPrice(entry: ChannelModelPricing, row: PricingRow): ChannelModelPricing {
  return {
    ...entry,
    platform: row.platform || 'openai',
    billing_mode: row.billingMode,
    input_price: parseTokenPrice(row.inputPriceText),
    output_price: parseTokenPrice(row.outputPriceText),
    cache_read_price: parseTokenPrice(row.cacheReadPriceText),
    cache_write_price: parseTokenPrice(row.cacheWritePriceText),
    per_request_price: parseDirectPrice(row.perRequestPriceText),
    public_visible: row.publicVisible,
    api_enabled: row.apiEnabled,
  }
}

async function saveRow(row: PricingRow): Promise<void> {
  row.saving = true
  row.error = ''
  row.saved = false
  try {
    validateRow(row)
    const channel = row.channelId ? channels.value.find(item => item.id === row.channelId) : primaryChannel.value
    if (!channel) throw new Error('没有可写入的渠道，请先创建并启用上游渠道')

    const pricing = (channel.model_pricing || []).map(clonePricing)
    if (row.source === 'channel' && row.entryIndex != null && pricing[row.entryIndex]) {
      const entry = pricing[row.entryIndex]
      const remainingModels = (entry.models || []).filter(model => canonicalPublicModelName(model).toLowerCase() !== row.model.toLowerCase())
      if ((entry.models || []).length > 1) {
        entry.models = remainingModels
        const newEntry = applyRowPrice({ ...entry, id: undefined, models: [row.model], intervals: [] }, row)
        pricing.push(newEntry)
      } else {
        pricing[row.entryIndex] = applyRowPrice({ ...entry, models: [row.model] }, row)
      }
    } else {
      pricing.push(applyRowPrice({
        platform: row.platform || 'openai',
        models: [row.model],
        billing_mode: row.billingMode,
        input_price: null,
        output_price: null,
        cache_write_price: null,
        cache_read_price: null,
        image_output_price: null,
        per_request_price: null,
        intervals: [],
      }, row))
    }

    const cleanedPricing = pricing.filter(entry => (entry.models || []).length > 0)
    await channelsAPI.update(channel.id, { model_pricing: cleanedPricing })
    row.dirty = false
    row.saved = true
    await load()
  } catch (err: any) {
    row.error = err?.message || '保存失败'
  } finally {
    row.saving = false
  }
}

async function saveDirtyRows(): Promise<void> {
  bulkSaving.value = true
  try {
    for (const row of [...dirtyRows.value]) {
      await saveRow(row)
    }
  } finally {
    bulkSaving.value = false
  }
}

function formatDate(value?: string): string {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

onMounted(load)
</script>
