<template>
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <div v-else class="min-h-screen bg-white text-slate-950 dark:bg-dark-950 dark:text-white">
    <PublicTopNav />

    <main>
      <section class="border-b border-primary-100 bg-gradient-to-br from-primary-50 via-white to-cyan-50 dark:border-dark-800 dark:from-dark-950 dark:via-dark-950 dark:to-primary-950/30">
        <div class="mx-auto grid max-w-7xl gap-10 px-4 py-16 sm:px-6 lg:grid-cols-[1fr_0.9fr] lg:px-8 lg:py-20">
          <div class="flex flex-col justify-center">
            <p class="mb-4 text-sm font-semibold uppercase tracking-wide text-primary-700 dark:text-primary-300">
              {{ siteName }}
            </p>
            <h1 class="max-w-3xl text-4xl font-bold leading-tight text-slate-950 dark:text-white sm:text-5xl lg:text-6xl">
              {{ copy.heroTitle }}
            </h1>
            <p class="mt-6 max-w-2xl text-lg leading-8 text-slate-600 dark:text-dark-200">
              {{ copy.heroDescription }}
            </p>

            <div v-if="isAuthenticated" class="mt-6 rounded-lg border border-primary-100 bg-white/80 p-4 text-sm text-slate-700 shadow-sm dark:border-primary-500/20 dark:bg-dark-900/70 dark:text-dark-200">
              <div class="font-semibold text-slate-950 dark:text-white">{{ copy.signedInAs }}</div>
              <div class="mt-1">{{ userIdentity }}</div>
            </div>

            <div class="mt-8 flex flex-col gap-3 sm:flex-row">
              <router-link
                :to="isAuthenticated ? '/dashboard' : '/register'"
                class="inline-flex items-center justify-center rounded-md bg-primary-600 px-5 py-3 text-sm font-bold text-white shadow-sm shadow-primary-600/20 transition-colors hover:bg-primary-700"
              >
                {{ isAuthenticated ? copy.goDashboard : copy.register }}
              </router-link>
              <router-link
                v-if="isAuthenticated"
                to="/purchase"
                class="inline-flex items-center justify-center rounded-md border border-primary-500 px-5 py-3 text-sm font-bold text-primary-700 transition-colors hover:bg-primary-50 dark:text-primary-300 dark:hover:bg-primary-500/10"
              >
                {{ copy.recharge }}
              </router-link>
              <router-link
                to="/models"
                class="inline-flex items-center justify-center rounded-md border border-primary-500 px-5 py-3 text-sm font-bold text-primary-700 transition-colors hover:bg-primary-50 dark:text-primary-300 dark:hover:bg-primary-500/10"
              >
                {{ copy.models }}
              </router-link>
              <router-link
                to="/docs"
                class="inline-flex items-center justify-center rounded-md border border-slate-300 px-5 py-3 text-sm font-bold text-slate-700 transition-colors hover:bg-white dark:border-dark-700 dark:text-dark-100 dark:hover:bg-dark-900"
              >
                {{ copy.docs }}
              </router-link>
            </div>

            <div class="mt-10 grid max-w-xl grid-cols-3 gap-5 border-t border-primary-200 pt-8 dark:border-dark-800">
              <div>
                <div class="text-3xl font-bold text-primary-700 dark:text-primary-300">{{ models.length }}</div>
                <div class="mt-1 text-sm text-slate-500 dark:text-dark-300">{{ copy.statsModels }}</div>
              </div>
              <div>
                <div class="text-3xl font-bold text-primary-700 dark:text-primary-300">1</div>
                <div class="mt-1 text-sm text-slate-500 dark:text-dark-300">{{ copy.statsGateway }}</div>
              </div>
              <div>
                <div class="text-3xl font-bold text-primary-700 dark:text-primary-300">{{ copy.statsRealtime }}</div>
                <div class="mt-1 text-sm text-slate-500 dark:text-dark-300">{{ copy.statsBilling }}</div>
              </div>
            </div>
          </div>

          <div class="flex min-w-0 items-center">
            <div class="w-full min-w-0 rounded-lg bg-gradient-to-br from-primary-600 via-primary-500 to-cyan-500 p-6 text-white shadow-xl shadow-primary-900/10 sm:p-8">
              <div class="flex items-center gap-3">
                <Icon name="terminal" size="xl" />
                <div>
                  <p class="text-sm font-semibold text-primary-50">{{ copy.apiExample }}</p>
                  <h2 class="text-2xl font-bold">{{ copy.openaiCompatible }}</h2>
                </div>
              </div>
              <pre class="mt-8 max-w-full overflow-x-auto rounded-md bg-primary-950/35 p-4 text-sm leading-7 text-primary-50 sm:p-5"><code>curl {{ apiBaseUrl }}/chat/completions \
  -H "Authorization: Bearer sk-your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-v4-flash",
    "messages": [{"role":"user","content":"Hello"}]
  }'</code></pre>
              <div class="mt-6 grid gap-3 sm:grid-cols-3">
                <div class="rounded-md bg-white/15 p-3">
                  <div class="text-sm text-primary-50">{{ copy.auth }}</div>
                  <div class="mt-1 font-semibold">Bearer Token</div>
                </div>
                <div class="rounded-md bg-white/15 p-3">
                  <div class="text-sm text-primary-50">{{ copy.protocol }}</div>
                  <div class="mt-1 font-semibold">OpenAI API</div>
                </div>
                <div class="rounded-md bg-white/15 p-3">
                  <div class="text-sm text-primary-50">{{ copy.billing }}</div>
                  <div class="mt-1 font-semibold">{{ copy.usageBased }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
        <div class="mb-8 flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
          <div>
            <p class="text-sm font-semibold text-primary-700 dark:text-primary-300">{{ copy.modelMarketplace }}</p>
            <h2 class="mt-2 text-3xl font-bold text-slate-950 dark:text-white">{{ copy.featuredModels }}</h2>
          </div>
          <router-link to="/models" class="text-sm font-semibold text-primary-700 hover:text-primary-800 dark:text-primary-300">
            {{ copy.viewAllModels }}
          </router-link>
        </div>
        <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
          <router-link
            v-for="model in models"
            :key="model.id"
            :to="{ path: '/models', query: { q: model.displayName } }"
            class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm transition hover:border-primary-200 hover:shadow-card-hover dark:border-dark-800 dark:bg-dark-900"
          >
            <div class="flex items-center justify-between gap-3">
              <h3 class="truncate text-lg font-bold text-slate-950 dark:text-white">{{ model.displayName }}</h3>
              <span class="rounded-full bg-primary-50 px-2.5 py-1 text-xs font-semibold text-primary-700 dark:bg-primary-500/10 dark:text-primary-300">
                {{ model.provider }}
              </span>
            </div>
            <p class="mt-3 line-clamp-3 text-sm leading-6 text-slate-600 dark:text-dark-300">
              {{ modelDescription(model) }}
            </p>
            <div class="mt-4 text-sm text-slate-500 dark:text-dark-400">
              {{ modelBillingLabel(model) }} / {{ model.unit }}
            </div>
          </router-link>
        </div>
      </section>

      <section class="border-y border-slate-200 bg-slate-50 dark:border-dark-800 dark:bg-dark-900/60">
        <div class="mx-auto grid max-w-7xl gap-6 px-4 py-16 sm:px-6 lg:grid-cols-3 lg:px-8">
          <div v-for="item in featureCards" :key="item.title" class="rounded-lg border border-slate-200 bg-white p-6 dark:border-dark-800 dark:bg-dark-950">
            <Icon :name="item.icon" size="lg" class="text-primary-600" />
            <h3 class="mt-4 text-xl font-bold">{{ item.title }}</h3>
            <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-dark-300">
              {{ item.description }}
            </p>
          </div>
        </div>
      </section>
    </main>

    <footer class="border-t border-slate-200 bg-slate-950 text-slate-300 dark:border-dark-800">
      <div class="mx-auto grid max-w-7xl gap-8 px-4 py-10 sm:px-6 lg:grid-cols-[1.2fr_0.8fr_0.8fr] lg:px-8">
        <div>
          <p class="font-bold text-white">Ekkamai Technology (Hong Kong) Company Limited</p>
          <address class="mt-3 max-w-xl text-sm not-italic leading-6 text-slate-400">
            <span class="block">Room 701, Unit 127, 7/F, Tower B, New Mandarin Plaza, 14 Science Museum Road, Tsim Sha Tsui, Kowloon</span>
            <span class="mt-1 block">Phone: 66841850843</span>
          </address>
          <a :href="contactHref" class="mt-4 inline-flex text-sm font-semibold text-cyan-300 hover:text-cyan-200">{{ contactText }}</a>
        </div>

        <div>
          <p class="text-sm font-bold uppercase tracking-wide text-white">{{ copy.serviceInformation }}</p>
          <div class="mt-3 grid gap-2 text-sm text-slate-400">
            <router-link to="/models" class="hover:text-white">{{ copy.modelsAndPricing }}</router-link>
            <router-link to="/legal/company" class="hover:text-white">{{ copy.companyInformation }}</router-link>
            <router-link to="/legal/payment-process" class="hover:text-white">{{ copy.paymentProcess }}</router-link>
            <router-link to="/support" class="hover:text-white">{{ copy.customerSupport }}</router-link>
          </div>
        </div>

        <div>
          <p class="text-sm font-bold uppercase tracking-wide text-white">{{ copy.policies }}</p>
          <div class="mt-3 grid gap-2 text-sm text-slate-400">
            <router-link to="/legal/terms" class="hover:text-white">{{ copy.terms }}</router-link>
            <router-link to="/legal/privacy" class="hover:text-white">{{ copy.privacy }}</router-link>
            <router-link to="/legal/delivery-policy" class="hover:text-white">{{ copy.delivery }}</router-link>
            <router-link to="/legal/refund-policy" class="hover:text-white">{{ copy.refunds }}</router-link>
          </div>
        </div>
      </div>
      <div class="border-t border-white/10">
        <div class="mx-auto max-w-7xl px-4 py-5 text-xs text-slate-500 sm:px-6 lg:px-8">
          © 2026 Ekkamai Technology (Hong Kong) Company Limited. All rights reserved.
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import PublicTopNav from '@/components/public/PublicTopNav.vue'
import Icon from '@/components/icons/Icon.vue'
import {
  dedupePublicModels,
  localizePublicModelDescription,
  publicModels,
  type PublicModelInfo,
} from '@/constants/publicModels'

const CONTACT_EMAIL = 'service@tokenapifuel.com'

const { locale } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.siteName || 'OneAPI')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const userIdentity = computed(() => authStore.user?.email || authStore.user?.username || 'OneAPI User')
const models = dedupePublicModels(publicModels).filter(model => model.endpointTypes.includes('openai:/v1/chat/completions')).slice(0, 5)
const currentLocale = computed(() => String(locale.value || 'zh-CN'))

const apiBaseUrl = computed(() => {
  const base = appStore.apiBaseUrl?.trim() || 'https://tokenapifuel.com'
  return `${base.replace(/^http:\/\/tokenapifuel\.com/i, 'https://tokenapifuel.com').replace(/\/$/, '')}/v1`
})

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const localized = {
  zh: {
    heroTitle: '一个 API Key，调用主流 AI 模型',
    heroDescription: '为开发者和运营团队提供统一网关、模型价格、API Key、充值余额扣费和用量记录。注册后即可创建密钥，按模型调用并实时查看消耗。',
    register: '免费注册',
    goDashboard: '进入控制台',
    recharge: '充值',
    models: '模型广场',
    docs: '开发文档',
    signedInAs: '当前已登录',
    statsModels: '主推模型',
    statsGateway: '统一网关',
    statsRealtime: '实时',
    statsBilling: '用量扣费',
    apiExample: 'API 配置示例',
    openaiCompatible: 'OpenAI 兼容调用',
    auth: '鉴权',
    protocol: '协议',
    billing: '计费',
    usageBased: '按用量扣费',
    featuredModels: '主推模型',
    modelMarketplace: '模型广场',
    tokenBilling: '按量计费',
    requestBilling: '按次计费',
    viewAllModels: '查看全部模型',
    apiDocs: 'API 接入说明',
    terms: '服务条款',
    serviceInformation: '服务信息',
    modelsAndPricing: '产品与价格',
    companyInformation: '公司信息',
    paymentProcess: '支付流程',
    customerSupport: '客户服务',
    policies: '协议与政策',
    privacy: '隐私政策',
    delivery: '数字交付与物流政策',
    refunds: '退换货及退款政策',
    featureCards: [
      ['key', '创建 API Key', '用户在控制台创建密钥，密钥额度与账户充值余额关联，方便运营和交付测试。'],
      ['creditCard', '充值余额', '钱包按固定泰铢金额充值，调用不同模型时按 Token 和模型价格实时扣费。'],
      ['chart', '用量可查', '每次请求记录模型、Token usage、扣费和调用状态，便于客户验收。'],
    ],
  },
  en: {
    heroTitle: 'One API key for mainstream AI models',
    heroDescription: 'A unified gateway for API keys, model pricing, THB balance billing, and usage records. Create a key after signup and track every request in real time.',
    register: 'Sign up',
    goDashboard: 'Open console',
    recharge: 'Top up',
    models: 'Models',
    docs: 'Docs',
    signedInAs: 'Signed in as',
    statsModels: 'Featured models',
    statsGateway: 'Unified gateway',
    statsRealtime: 'Real-time',
    statsBilling: 'billing',
    apiExample: 'API example',
    openaiCompatible: 'OpenAI-compatible',
    auth: 'Auth',
    protocol: 'Protocol',
    billing: 'Billing',
    usageBased: 'Usage based',
    featuredModels: 'Featured models',
    modelMarketplace: 'Model Marketplace',
    tokenBilling: 'Token billing',
    requestBilling: 'Per request',
    viewAllModels: 'View all models',
    apiDocs: 'API guide',
    terms: 'Terms',
    serviceInformation: 'Service information',
    modelsAndPricing: 'Products & pricing',
    companyInformation: 'Company information',
    paymentProcess: 'Payment process',
    customerSupport: 'Customer support',
    policies: 'Legal & policies',
    privacy: 'Privacy policy',
    delivery: 'Digital delivery policy',
    refunds: 'Refund policy',
    featureCards: [
      ['key', 'Create API keys', 'Create keys in the console. Key access is tied to account balance for easy operations and delivery testing.'],
      ['creditCard', 'THB top-up', 'Top up your wallet with fixed THB amounts. Calls are charged by model price and token usage.'],
      ['chart', 'Usage records', 'Every request records model, token usage, billing amount, and status for customer acceptance.'],
    ],
  },
  th: {
    heroTitle: 'API Key เดียว เรียกใช้โมเดล AI หลักได้',
    heroDescription: 'เกตเวย์รวมสำหรับ API Key ราคาโมเดล ยอดเงิน THB การหักเงินตามการใช้งาน และบันทึกการเรียกใช้แบบเรียลไทม์',
    register: 'สมัครใช้งาน',
    goDashboard: 'เปิดคอนโซล',
    recharge: 'เติมเงิน',
    models: 'โมเดล',
    docs: 'เอกสาร',
    signedInAs: 'เข้าสู่ระบบเป็น',
    statsModels: 'โมเดลแนะนำ',
    statsGateway: 'เกตเวย์เดียว',
    statsRealtime: 'เรียลไทม์',
    statsBilling: 'คิดค่าบริการ',
    apiExample: 'ตัวอย่าง API',
    openaiCompatible: 'เข้ากันได้กับ OpenAI',
    auth: 'ยืนยันตัวตน',
    protocol: 'โปรโตคอล',
    billing: 'การคิดเงิน',
    usageBased: 'คิดตามการใช้งาน',
    featuredModels: 'โมเดลแนะนำ',
    modelMarketplace: 'ตลาดโมเดล',
    tokenBilling: 'คิดเงินตาม Token',
    requestBilling: 'คิดต่อครั้ง',
    viewAllModels: 'ดูโมเดลทั้งหมด',
    apiDocs: 'คู่มือ API',
    terms: 'เงื่อนไขบริการ',
    serviceInformation: 'ข้อมูลบริการ',
    modelsAndPricing: 'ผลิตภัณฑ์และราคา',
    companyInformation: 'ข้อมูลบริษัท',
    paymentProcess: 'ขั้นตอนการชำระเงิน',
    customerSupport: 'ฝ่ายบริการลูกค้า',
    policies: 'ข้อกำหนดและนโยบาย',
    privacy: 'นโยบายความเป็นส่วนตัว',
    delivery: 'นโยบายการส่งมอบดิจิทัล',
    refunds: 'นโยบายคืนเงิน',
    featureCards: [
      ['key', 'สร้าง API Key', 'สร้างคีย์ในคอนโซลและผูกกับยอดเงินในบัญชีเพื่อทดสอบและส่งมอบได้ง่าย'],
      ['creditCard', 'เติมเงิน THB', 'เติมเงินเป็นจำนวนเงินบาทแบบคงที่ และหักตามราคาโมเดลกับ Token ที่ใช้จริง'],
      ['chart', 'ตรวจสอบการใช้งาน', 'บันทึกโมเดล Token usage จำนวนเงินที่หัก และสถานะของทุกคำขอ'],
    ],
  },
}

const copy = computed(() => {
  const current = String(locale.value || '').toLowerCase()
  if (current.startsWith('th')) return localized.th
  if (current.startsWith('en')) return localized.en
  return localized.zh
})

type FeatureIcon = 'key' | 'creditCard' | 'chart'

const featureCards = computed(() =>
  copy.value.featureCards.map(([icon, title, description]) => ({ icon: icon as FeatureIcon, title, description }))
)

function modelDescription(model: PublicModelInfo): string {
  return localizePublicModelDescription(model, currentLocale.value)
}

function modelBillingLabel(model: PublicModelInfo): string {
  return model.billingMode === 'request' ? copy.value.requestBilling : copy.value.tokenBilling
}

const contactText = computed(() => appStore.contactInfo || CONTACT_EMAIL)
const contactHref = computed(() => `mailto:${CONTACT_EMAIL}`)
</script>
