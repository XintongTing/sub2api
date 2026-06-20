<template>
  <div class="min-h-screen bg-white text-slate-950 dark:bg-dark-950 dark:text-white">
    <PublicTopNav />

    <main class="grid min-h-[calc(100vh-4rem)] lg:grid-cols-[300px_1fr]">
      <aside class="border-b border-slate-200 bg-white p-4 dark:border-dark-800 dark:bg-dark-950 lg:border-b-0 lg:border-r">
        <div class="mb-4 flex items-center justify-between">
          <h1 class="text-lg font-bold">{{ text.filters }}</h1>
          <button
            type="button"
            class="rounded-md border border-slate-200 px-3 py-1.5 text-sm font-semibold text-slate-600 transition hover:border-primary-200 hover:bg-primary-50 hover:text-primary-700 dark:border-dark-800 dark:text-dark-200 dark:hover:bg-primary-500/10"
            @click="resetFilters"
          >
            {{ text.reset }}
          </button>
        </div>

        <FilterSection :title="text.provider">
          <button
            v-for="provider in providerOptions"
            :key="provider.value"
            type="button"
            class="filter-pill"
            :class="{ 'filter-pill-active': selectedProvider === provider.value }"
            @click="selectedProvider = provider.value"
          >
            {{ provider.label }}
          </button>
        </FilterSection>

        <FilterSection :title="text.billingType">
          <button
            v-for="mode in billingModes"
            :key="mode.value"
            type="button"
            class="filter-pill"
            :class="{ 'filter-pill-active': selectedBilling === mode.value }"
            @click="selectedBilling = mode.value"
          >
            {{ mode.label }}
          </button>
        </FilterSection>

        <FilterSection :title="text.tags">
          <button
            v-for="tag in tagOptions"
            :key="tag.value"
            type="button"
            class="filter-pill"
            :class="{ 'filter-pill-active': selectedTag === tag.value }"
            @click="selectedTag = tag.value"
          >
            {{ tag.label }}
          </button>
        </FilterSection>

        <FilterSection :title="text.endpointType">
          <button
            v-for="endpoint in endpointTypes"
            :key="endpoint.value"
            type="button"
            class="filter-pill"
            :class="{ 'filter-pill-active': selectedEndpoint === endpoint.value }"
            @click="selectedEndpoint = endpoint.value"
          >
            {{ endpoint.label }}
          </button>
        </FilterSection>
      </aside>

      <section class="min-w-0 bg-slate-50/80 dark:bg-dark-950">
        <div class="border-b border-slate-200 bg-white p-4 dark:border-dark-800 dark:bg-dark-950 lg:p-6">
          <div class="rounded-lg bg-gradient-to-r from-primary-700 via-primary-600 to-cyan-500 p-6 text-white shadow-sm">
            <div class="flex flex-col justify-between gap-4 xl:flex-row xl:items-end">
              <div>
                <div class="flex flex-wrap items-center gap-3">
                  <h2 class="text-2xl font-bold">{{ text.allModels }}</h2>
                  <span class="rounded-full bg-white px-3 py-1 text-sm font-semibold text-primary-700">
                    {{ text.modelCount(filteredModels.length) }}
                  </span>
                </div>
                <p class="mt-2 max-w-3xl text-sm leading-6 text-primary-50">
                  {{ text.heroDescription }}
                </p>
              </div>
              <router-link
                to="/docs"
                class="inline-flex items-center justify-center rounded-md bg-white px-4 py-2 text-sm font-bold text-primary-700 transition hover:bg-primary-50"
              >
                {{ text.viewDocs }}
              </router-link>
            </div>
          </div>

          <div class="mt-4 flex flex-col gap-3 rounded-lg border border-slate-200 bg-white p-3 dark:border-dark-800 dark:bg-dark-900 xl:flex-row xl:items-center">
            <div class="relative min-w-0 flex-1">
              <Icon name="search" size="sm" class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
              <input
                v-model="search"
                type="search"
                class="w-full rounded-md border border-slate-200 bg-slate-50 py-2 pl-9 pr-3 text-sm text-slate-900 outline-none transition focus:border-primary-400 focus:ring-2 focus:ring-primary-500/20 dark:border-dark-800 dark:bg-dark-950 dark:text-white"
                :placeholder="text.searchPlaceholder"
              />
            </div>
            <div class="flex flex-wrap items-center gap-3">
              <button type="button" class="toolbar-button" @click="copyBaseUrl">
                <Icon name="copy" size="sm" />
                {{ text.copyBaseUrl }}
              </button>
              <label class="flex items-center gap-2 text-sm font-medium text-slate-600 dark:text-dark-200">
                {{ text.showPrices }}
                <input v-model="showPrices" type="checkbox" class="h-4 w-4 accent-primary-600" />
              </label>
              <button
                type="button"
                class="toolbar-button"
                :class="{ 'toolbar-button-active': viewMode === 'card' }"
                @click="viewMode = 'card'"
              >
                <Icon name="grid" size="sm" />
                {{ text.cardView }}
              </button>
              <button
                type="button"
                class="toolbar-button"
                :class="{ 'toolbar-button-active': viewMode === 'table' }"
                @click="viewMode = 'table'"
              >
                <Icon name="document" size="sm" />
                {{ text.tableView }}
              </button>
            </div>
          </div>
        </div>

        <div class="p-4 lg:p-6">
          <div v-if="isLoading" class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <div v-for="i in 6" :key="i" class="h-52 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-dark-800 dark:bg-dark-900"></div>
          </div>

          <div v-else-if="viewMode === 'card'" class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <article
              v-for="model in filteredModels"
              :key="model.name"
              class="model-card"
              @click="selectedModel = model"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0">
                  <h3 class="truncate text-xl font-bold text-slate-950 dark:text-white">{{ model.displayName }}</h3>
                  <p class="mt-1 text-sm text-slate-500 dark:text-dark-400">{{ model.provider }}</p>
                </div>
                <button type="button" class="rounded-md p-1.5 text-slate-400 hover:bg-primary-50 hover:text-primary-700 dark:hover:bg-primary-500/10" @click.stop="copyModelName(model.name)">
                  <Icon name="copy" size="sm" />
                </button>
              </div>

              <div class="mt-3 space-y-1.5 text-sm text-slate-700 dark:text-dark-200">
                <PriceLine :label="text.inputPrice" :value="formatTokenPrice(model.inputPrice)" :visible="showPrices" :hidden-label="text.hidden" />
                <PriceLine :label="text.outputPrice" :value="formatTokenPrice(model.outputPrice)" :visible="showPrices" :hidden-label="text.hidden" />
                <PriceLine v-if="model.cacheReadPrice" :label="text.cacheReadPrice" :value="formatTokenPrice(model.cacheReadPrice)" :visible="showPrices" :hidden-label="text.hidden" />
                <PriceLine v-if="model.cacheWritePrice" :label="text.cacheWritePrice" :value="formatTokenPrice(model.cacheWritePrice)" :visible="showPrices" :hidden-label="text.hidden" />
                <PriceLine v-if="model.perRequestPrice" :label="text.perRequestPrice" :value="formatRequestPrice(model.perRequestPrice)" :visible="showPrices" :hidden-label="text.hidden" />
              </div>

              <p class="mt-4 line-clamp-3 min-h-[4.5rem] text-sm leading-6 text-slate-600 dark:text-dark-300">
                {{ model.description }}
              </p>
              <div class="mt-4 flex flex-wrap gap-2">
                <span v-for="tag in model.tags.slice(0, 4)" :key="tag" class="tag-pill">{{ tag }}</span>
                <span v-if="model.tags.length > 4" class="tag-pill">+{{ model.tags.length - 4 }}</span>
              </div>
            </article>
          </div>

          <div v-else class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-dark-800 dark:bg-dark-900">
            <table class="min-w-full divide-y divide-slate-200 dark:divide-dark-800">
              <thead class="bg-slate-50 text-left text-sm font-semibold text-slate-600 dark:bg-dark-800 dark:text-dark-200">
                <tr>
                  <th class="px-4 py-3">{{ text.model }}</th>
                  <th class="px-4 py-3">{{ text.provider }}</th>
                  <th class="px-4 py-3">{{ text.billingType }}</th>
                  <th class="px-4 py-3">{{ text.inputPrice }}</th>
                  <th class="px-4 py-3">{{ text.outputPrice }}</th>
                  <th class="px-4 py-3">{{ text.endpointType }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 text-sm dark:divide-dark-800">
                <tr v-for="model in filteredModels" :key="model.name" class="cursor-pointer hover:bg-primary-50/60 dark:hover:bg-primary-500/10" @click="selectedModel = model">
                  <td class="px-4 py-3 font-semibold">{{ model.displayName }}</td>
                  <td class="px-4 py-3">{{ model.provider }}</td>
                  <td class="px-4 py-3">{{ billingLabel(model.billingMode) }}</td>
                  <td class="px-4 py-3">{{ showPrices ? formatTokenPrice(model.inputPrice) : text.hidden }}</td>
                  <td class="px-4 py-3">{{ showPrices ? formatTokenPrice(model.outputPrice) : text.hidden }}</td>
                  <td class="px-4 py-3">{{ model.endpointTypes[0] || '/v1/chat/completions' }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-if="!isLoading && filteredModels.length === 0" class="rounded-lg border border-slate-200 bg-white p-12 text-center dark:border-dark-800 dark:bg-dark-900">
            <Icon name="search" size="xl" class="mx-auto text-slate-300" />
            <h3 class="mt-4 text-lg font-bold">{{ text.emptyTitle }}</h3>
            <p class="mt-2 text-sm text-slate-500">{{ text.emptyDescription }}</p>
          </div>

          <p v-if="loadError" class="mt-4 rounded-md border border-amber-200 bg-amber-50 p-3 text-sm text-amber-800 dark:border-amber-900/50 dark:bg-amber-950/20 dark:text-amber-200">
            {{ loadError }}{{ text.fallbackNotice }}
          </p>
        </div>
      </section>
    </main>

    <div v-if="selectedModel" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 p-4" @click.self="selectedModel = null">
      <section class="max-h-[90vh] w-full max-w-3xl overflow-y-auto rounded-lg bg-white shadow-xl dark:bg-dark-900">
        <div class="flex items-start justify-between gap-4 border-b border-slate-200 p-5 dark:border-dark-800">
          <div>
            <h2 class="text-2xl font-bold">{{ selectedModel.displayName }}</h2>
            <p class="mt-1 text-sm text-slate-500 dark:text-dark-300">{{ selectedModel.provider }} / {{ billingLabel(selectedModel.billingMode) }}</p>
          </div>
          <button type="button" class="rounded-md p-2 text-slate-500 hover:bg-slate-100 dark:hover:bg-dark-800" @click="selectedModel = null">
            <Icon name="x" size="md" />
          </button>
        </div>
        <div class="space-y-6 p-5">
          <section>
            <h3 class="font-bold">{{ text.basicInfo }}</h3>
            <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-dark-300">{{ selectedModel.description }}</p>
          </section>
          <section>
            <h3 class="font-bold">{{ text.apiEndpoint }}</h3>
            <div class="mt-2 rounded-md border border-slate-200 bg-slate-50 p-3 text-sm dark:border-dark-800 dark:bg-dark-950">
              <div class="font-mono">POST {{ apiBaseUrl }}/chat/completions</div>
              <div class="mt-1 text-slate-500">OpenAI compatible / model = {{ selectedModel.name }}</div>
            </div>
          </section>
          <section>
            <h3 class="font-bold">{{ text.priceSummary }}</h3>
            <div class="mt-3 grid gap-3 sm:grid-cols-2">
              <div class="price-card"><span>{{ text.inputPrice }}</span><strong>{{ formatTokenPrice(selectedModel.inputPrice) }}</strong></div>
              <div class="price-card"><span>{{ text.outputPrice }}</span><strong>{{ formatTokenPrice(selectedModel.outputPrice) }}</strong></div>
              <div class="price-card"><span>{{ text.cacheReadPrice }}</span><strong>{{ formatTokenPrice(selectedModel.cacheReadPrice) }}</strong></div>
              <div class="price-card"><span>{{ text.cacheWritePrice }}</span><strong>{{ formatTokenPrice(selectedModel.cacheWritePrice) }}</strong></div>
            </div>
          </section>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import PublicTopNav from '@/components/public/PublicTopNav.vue'
import Icon from '@/components/icons/Icon.vue'
import { listPublicModels, type PublicModelDTO } from '@/api/publicModels'
import { publicModels, type PublicModelInfo } from '@/constants/publicModels'
import { useAppStore } from '@/stores'

const COPY = {
  'zh-CN': {
    filters: '筛选',
    reset: '重置',
    provider: '供应商',
    billingType: '计费类型',
    tags: '标签',
    endpointType: '端点类型',
    all: '全部',
    allProviders: '全部供应商',
    allModels: '全部模型',
    modelCount: (count: number) => `共 ${count} 个模型`,
    heroDescription: '查看本站已接入的 AI 模型与公开价格。实际扣费以后台模型价格、API Key 权限和账户充值余额为准。',
    viewDocs: '查看接入文档',
    searchPlaceholder: '模糊搜索模型名称',
    copyBaseUrl: '复制 Base URL',
    showPrices: '显示价格',
    cardView: '卡片视图',
    tableView: '表格视图',
    model: '模型',
    inputPrice: '输入价格',
    outputPrice: '补全价格',
    cacheReadPrice: '缓存读取价格',
    cacheWritePrice: '缓存创建价格',
    perRequestPrice: '按次价格',
    tokenBilling: '按量计费',
    requestBilling: '按次计费',
    hidden: '隐藏',
    unset: '后台未配置',
    each: '次',
    emptyTitle: '没有找到匹配模型',
    emptyDescription: '换一个关键词或重置筛选条件。',
    fallbackNotice: '，当前展示静态兜底模型。',
    basicInfo: '基本信息',
    apiEndpoint: 'API 端点',
    priceSummary: '价格摘要',
  },
  'zh-TW': {
    filters: '篩選',
    reset: '重置',
    provider: '供應商',
    billingType: '計費類型',
    tags: '標籤',
    endpointType: '端點類型',
    all: '全部',
    allProviders: '全部供應商',
    allModels: '全部模型',
    modelCount: (count: number) => `共 ${count} 個模型`,
    heroDescription: '查看本站已接入的 AI 模型與公開價格。實際扣費以後台模型價格、API Key 權限和帳戶充值餘額為準。',
    viewDocs: '查看接入文件',
    searchPlaceholder: '模糊搜尋模型名稱',
    copyBaseUrl: '複製 Base URL',
    showPrices: '顯示價格',
    cardView: '卡片視圖',
    tableView: '表格視圖',
    model: '模型',
    inputPrice: '輸入價格',
    outputPrice: '補全價格',
    cacheReadPrice: '快取讀取價格',
    cacheWritePrice: '快取建立價格',
    perRequestPrice: '按次價格',
    tokenBilling: '按量計費',
    requestBilling: '按次計費',
    hidden: '隱藏',
    unset: '後台未配置',
    each: '次',
    emptyTitle: '沒有找到匹配模型',
    emptyDescription: '換一個關鍵字或重置篩選條件。',
    fallbackNotice: '，目前展示靜態兜底模型。',
    basicInfo: '基本資訊',
    apiEndpoint: 'API 端點',
    priceSummary: '價格摘要',
  },
  en: {
    filters: 'Filters',
    reset: 'Reset',
    provider: 'Provider',
    billingType: 'Billing',
    tags: 'Tags',
    endpointType: 'Endpoint',
    all: 'All',
    allProviders: 'All providers',
    allModels: 'All models',
    modelCount: (count: number) => `${count} models`,
    heroDescription: 'Explore the AI models enabled on this gateway and their public prices. Actual billing follows backend model pricing, API key permissions, and account top-up balance.',
    viewDocs: 'View docs',
    searchPlaceholder: 'Search model names',
    copyBaseUrl: 'Copy Base URL',
    showPrices: 'Show prices',
    cardView: 'Cards',
    tableView: 'Table',
    model: 'Model',
    inputPrice: 'Input price',
    outputPrice: 'Output price',
    cacheReadPrice: 'Cache read',
    cacheWritePrice: 'Cache write',
    perRequestPrice: 'Per request',
    tokenBilling: 'Token billing',
    requestBilling: 'Per request',
    hidden: 'Hidden',
    unset: 'Not configured',
    each: 'request',
    emptyTitle: 'No models found',
    emptyDescription: 'Try another keyword or reset the filters.',
    fallbackNotice: '; showing the static fallback models.',
    basicInfo: 'Basic Info',
    apiEndpoint: 'API Endpoint',
    priceSummary: 'Price Summary',
  },
  th: {
    filters: 'ตัวกรอง',
    reset: 'รีเซ็ต',
    provider: 'ผู้ให้บริการ',
    billingType: 'การคิดเงิน',
    tags: 'แท็ก',
    endpointType: 'Endpoint',
    all: 'ทั้งหมด',
    allProviders: 'ผู้ให้บริการทั้งหมด',
    allModels: 'โมเดลทั้งหมด',
    modelCount: (count: number) => `${count} โมเดล`,
    heroDescription: 'ดูโมเดล AI ที่เปิดใช้งานบนเกตเวย์นี้และราคาสาธารณะ การหักเงินจริงอิงตามราคาหลังบ้าน สิทธิ์ API Key และยอดเติมเงินในบัญชี',
    viewDocs: 'ดูเอกสาร',
    searchPlaceholder: 'ค้นหาชื่อโมเดล',
    copyBaseUrl: 'คัดลอก Base URL',
    showPrices: 'แสดงราคา',
    cardView: 'การ์ด',
    tableView: 'ตาราง',
    model: 'โมเดล',
    inputPrice: 'ราคา Input',
    outputPrice: 'ราคา Output',
    cacheReadPrice: 'ราคา Cache Read',
    cacheWritePrice: 'ราคา Cache Write',
    perRequestPrice: 'ราคาต่อครั้ง',
    tokenBilling: 'คิดตาม Token',
    requestBilling: 'คิดต่อครั้ง',
    hidden: 'ซ่อน',
    unset: 'ยังไม่ได้ตั้งค่า',
    each: 'ครั้ง',
    emptyTitle: 'ไม่พบโมเดล',
    emptyDescription: 'ลองเปลี่ยนคำค้นหรือรีเซ็ตตัวกรอง',
    fallbackNotice: ' กำลังแสดงโมเดลสำรองแบบคงที่',
    basicInfo: 'ข้อมูลพื้นฐาน',
    apiEndpoint: 'API Endpoint',
    priceSummary: 'สรุปราคา',
  },
} as const

type CopyKey = keyof typeof COPY

const FilterSection = defineComponent({
  props: { title: { type: String, required: true } },
  setup(props, { slots }) {
    return () => h('section', { class: 'border-t border-slate-200 py-5 dark:border-dark-800' }, [
      h('h2', { class: 'mb-3 text-sm font-bold text-slate-950 dark:text-white' }, props.title),
      h('div', { class: 'grid grid-cols-2 gap-2' }, slots.default?.()),
    ])
  },
})

const PriceLine = defineComponent({
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true },
    visible: { type: Boolean, required: true },
    hiddenLabel: { type: String, required: true },
  },
  setup(props) {
    return () => h('div', { class: 'flex items-baseline gap-2' }, [
      h('span', { class: 'text-slate-500 dark:text-dark-400' }, props.label),
      h('strong', { class: 'font-semibold text-slate-950 dark:text-white' }, props.visible ? props.value : props.hiddenLabel),
    ])
  },
})

const route = useRoute()
const { locale } = useI18n()
const appStore = useAppStore()

const models = ref<PublicModelInfo[]>(publicModels)
const isLoading = ref(false)
const loadError = ref('')
const search = ref(typeof route.query.q === 'string' ? route.query.q : '')
const selectedProvider = ref('all')
const selectedBilling = ref('all')
const selectedTag = ref('all')
const selectedEndpoint = ref('all')
const viewMode = ref<'card' | 'table'>('card')
const showPrices = ref(true)
const selectedModel = ref<PublicModelInfo | null>(null)

const copyKey = computed<CopyKey>(() => {
  const value = String(locale.value || 'zh-CN')
  if (value.startsWith('zh-TW')) return 'zh-TW'
  if (value.startsWith('en')) return 'en'
  if (value.startsWith('th')) return 'th'
  return 'zh-CN'
})
const text = computed(() => COPY[copyKey.value])

const billingModes = computed(() => [
  { label: text.value.all, value: 'all' },
  { label: text.value.tokenBilling, value: 'token' },
  { label: text.value.requestBilling, value: 'request' },
])

const providerOptions = computed(() => [
  { label: text.value.allProviders, value: 'all' },
  ...Array.from(new Set(models.value.map(model => model.provider).filter(Boolean)))
    .sort()
    .map(provider => ({ label: provider, value: provider })),
])

const tagOptions = computed(() => [
  { label: text.value.all, value: 'all' },
  ...Array.from(new Set(models.value.flatMap(model => model.tags || []).filter(Boolean)))
    .sort()
    .map(tag => ({ label: tag, value: tag })),
])

const endpointTypes = computed(() => [
  { label: text.value.all, value: 'all' },
  ...Array.from(new Set(models.value.flatMap(model => model.endpointTypes || []).filter(Boolean)))
    .sort()
    .map(endpoint => ({ label: endpoint.replace(/^openai:/, ''), value: endpoint })),
])

const apiBaseUrl = computed(() => {
  const base = appStore.apiBaseUrl?.trim() || 'https://tokenapifuel.com'
  return `${base.replace(/^http:\/\/tokenapifuel\.com/i, 'https://tokenapifuel.com').replace(/\/$/, '')}/v1`
})

const filteredModels = computed(() => {
  const keyword = search.value.trim().toLowerCase()
  return models.value.filter((model) => {
    const matchesSearch = !keyword || [model.name, model.displayName, model.provider, model.description, ...model.tags]
      .filter(Boolean)
      .some(value => value.toLowerCase().includes(keyword))
    const matchesProvider = selectedProvider.value === 'all' || model.provider === selectedProvider.value
    const matchesBilling = selectedBilling.value === 'all' || model.billingMode === selectedBilling.value
    const matchesTag = selectedTag.value === 'all' || model.tags.includes(selectedTag.value)
    const matchesEndpoint = selectedEndpoint.value === 'all' || model.endpointTypes.includes(selectedEndpoint.value)
    return matchesSearch && matchesProvider && matchesBilling && matchesTag && matchesEndpoint
  })
})

watch(() => route.query.q, (value) => {
  search.value = typeof value === 'string' ? value : ''
})

onMounted(async () => {
  isLoading.value = true
  loadError.value = ''
  try {
    const remote = await listPublicModels()
    if (remote.length > 0) {
      models.value = remote.map(mapRemoteModel)
    }
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : 'Public model API is unavailable'
    models.value = publicModels
  } finally {
    isLoading.value = false
  }
})

function mapRemoteModel(model: PublicModelDTO): PublicModelInfo {
  const known = publicModels.find(item => item.name === model.name || item.displayName === model.name)
  return {
    id: model.name,
    name: model.name,
    displayName: known?.displayName || model.name,
    provider: model.provider || known?.provider || 'OneAPI',
    upstreamModel: model.name,
    type: known?.type || 'Chat',
    billing: model.billing_mode === 'request' ? text.value.requestBilling : text.value.tokenBilling,
    billingMode: model.billing_mode || known?.billingMode || 'token',
    currency: model.currency || known?.currency || 'THB',
    inputPrice: model.input_price ?? known?.inputPrice ?? null,
    outputPrice: model.output_price ?? known?.outputPrice ?? null,
    cacheReadPrice: model.cache_read_price ?? known?.cacheReadPrice ?? null,
    cacheWritePrice: model.cache_write_price ?? known?.cacheWritePrice ?? null,
    perRequestPrice: model.per_request_price ?? known?.perRequestPrice ?? null,
    unit: model.unit || known?.unit || '1M Tokens',
    endpointTypes: model.endpoint_types?.length ? model.endpoint_types : (known?.endpointTypes || ['openai:/v1/chat/completions']),
    description: model.description || known?.description || '',
    tags: model.tags?.length ? model.tags : (known?.tags || []),
  }
}

function resetFilters(): void {
  search.value = ''
  selectedProvider.value = 'all'
  selectedBilling.value = 'all'
  selectedTag.value = 'all'
  selectedEndpoint.value = 'all'
}

function billingLabel(mode: string): string {
  return mode === 'request' ? text.value.requestBilling : text.value.tokenBilling
}

function formatTokenPrice(value?: number | null): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) return text.value.unset
  return `฿${(value * 1_000_000).toFixed(4)} / 1M Tokens`
}

function formatRequestPrice(value?: number | null): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) return text.value.unset
  return `฿${value.toFixed(4)} / ${text.value.each}`
}

async function copyBaseUrl(): Promise<void> {
  await navigator.clipboard?.writeText(apiBaseUrl.value)
}

async function copyModelName(name: string): Promise<void> {
  await navigator.clipboard?.writeText(name)
}
</script>

<style scoped>
.filter-pill {
  @apply rounded-md border border-slate-200 bg-white px-3 py-2 text-sm font-semibold text-slate-600 transition hover:border-primary-200 hover:bg-primary-50 hover:text-primary-700 dark:border-dark-800 dark:bg-dark-900 dark:text-dark-200 dark:hover:bg-primary-500/10;
}

.filter-pill-active {
  @apply border-primary-200 bg-primary-50 text-primary-700 dark:border-primary-500/30 dark:bg-primary-500/10 dark:text-primary-300;
}

.toolbar-button {
  @apply inline-flex items-center gap-1.5 rounded-md border border-slate-200 bg-white px-3 py-2 text-sm font-semibold text-slate-600 transition hover:border-primary-200 hover:bg-primary-50 hover:text-primary-700 dark:border-dark-800 dark:bg-dark-900 dark:text-dark-200 dark:hover:bg-primary-500/10;
}

.toolbar-button-active {
  @apply border-primary-200 bg-primary-50 text-primary-700 dark:border-primary-500/30 dark:bg-primary-500/10 dark:text-primary-300;
}

.model-card {
  @apply cursor-pointer rounded-lg border border-slate-200 bg-white p-5 shadow-sm transition hover:border-primary-200 hover:shadow-card-hover dark:border-dark-800 dark:bg-dark-900;
}

.tag-pill {
  @apply rounded-full bg-primary-50 px-2.5 py-1 text-xs font-semibold text-primary-700 dark:bg-primary-500/10 dark:text-primary-300;
}

.price-card {
  @apply flex flex-col gap-1 rounded-md border border-slate-200 bg-slate-50 p-3 text-sm dark:border-dark-800 dark:bg-dark-950;
}

.price-card span {
  @apply text-slate-500 dark:text-dark-400;
}

.price-card strong {
  @apply text-slate-950 dark:text-white;
}
</style>
