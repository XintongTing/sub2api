<template>
  <div class="min-h-screen bg-white text-slate-950 dark:bg-dark-950 dark:text-white">
    <PublicTopNav />

    <main class="mx-auto grid max-w-7xl gap-8 px-4 py-8 sm:px-6 lg:grid-cols-[220px_1fr_220px] lg:px-8">
      <aside class="hidden lg:block">
        <nav class="sticky top-24 space-y-1 text-sm">
          <a
            v-for="item in docNav"
            :key="item.href"
            :href="item.href"
            class="block rounded-md px-3 py-2 font-semibold text-slate-600 hover:bg-primary-50 hover:text-primary-700 dark:text-dark-300 dark:hover:bg-primary-500/10 dark:hover:text-primary-300"
          >
            {{ item.label }}
          </a>
        </nav>
      </aside>

      <article class="min-w-0">
        <section class="rounded-lg border border-primary-100 bg-gradient-to-br from-primary-50 via-white to-cyan-50 p-6 dark:border-dark-800 dark:from-dark-900 dark:via-dark-900 dark:to-primary-950/30">
          <p class="text-sm font-semibold text-primary-700 dark:text-primary-300">{{ docs.kicker }}</p>
          <h1 class="mt-2 text-3xl font-bold text-slate-950 dark:text-white">{{ docs.title }}</h1>
          <p class="mt-4 max-w-3xl text-base leading-7 text-slate-600 dark:text-dark-200">
            {{ docs.description }}
          </p>
        </section>

        <section id="quickstart" class="doc-section">
          <h2>{{ docs.sections.quickstart }}</h2>
          <ol>
            <li v-for="step in docs.quickstartSteps" :key="step">{{ step }}</li>
          </ol>
        </section>

        <section id="base-url" class="doc-section">
          <h2>{{ sectionLabels.baseUrl }}</h2>
          <p>{{ docs.baseUrl.current }}</p>
          <pre><code>{{ gatewayBaseUrl }}</code></pre>
          <p>{{ docs.baseUrl.sdk }}</p>
          <pre><code>{{ apiBaseUrl }}</code></pre>
          <p>{{ docs.baseUrl.httpsNote }}</p>
        </section>

        <section id="auth" class="doc-section">
          <h2>{{ docs.sections.auth }}</h2>
          <p>{{ docs.authDescription }}</p>
          <pre><code>Authorization: Bearer sk-your-api-key
Content-Type: application/json</code></pre>
        </section>

        <section id="chat" class="doc-section">
          <h2>{{ sectionLabels.chat }}</h2>
          <p>{{ docs.chatDescription }}</p>
          <pre><code>POST /v1/chat/completions</code></pre>

          <h3>cURL</h3>
          <pre><code>{{ curlExample }}</code></pre>

          <h3>Python</h3>
          <pre><code>{{ pythonExample }}</code></pre>

          <h3>Node.js</h3>
          <pre><code>{{ nodeExample }}</code></pre>

          <h3>Go</h3>
          <pre><code>{{ goExample }}</code></pre>
        </section>

        <section id="codex" class="doc-section">
          <h2>{{ sectionLabels.codex }}</h2>
          <p>{{ docs.codexDescription }}</p>
          <p>{{ docs.codexConfigIntro }}</p>
          <pre><code>{{ codexConfigExample }}</code></pre>

          <p>{{ docs.authJsonExample }}</p>
          <pre><code>{
  "OPENAI_API_KEY": "sk-your-api-key"
}</code></pre>
        </section>

        <section id="models" class="doc-section">
          <h2>{{ docs.sections.models }}</h2>
          <p>{{ docs.modelsDescription }}</p>
          <div class="overflow-hidden rounded-lg border border-slate-200 dark:border-dark-800">
            <table class="min-w-full divide-y divide-slate-200 text-sm dark:divide-dark-800">
              <thead class="bg-slate-50 text-left font-semibold text-slate-600 dark:bg-dark-800 dark:text-dark-200">
                <tr>
                  <th class="px-4 py-3">{{ docs.table.model }}</th>
                  <th class="px-4 py-3">{{ docs.table.provider }}</th>
                  <th class="px-4 py-3">{{ docs.table.input }}</th>
                  <th class="px-4 py-3">{{ docs.table.output }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-dark-800">
                <tr v-for="model in models" :key="model.id">
                  <td class="px-4 py-3 font-semibold">{{ model.displayName }}</td>
                  <td class="px-4 py-3">{{ model.provider }}</td>
                  <td class="px-4 py-3">{{ formatTokenPrice(model.inputPrice) }}</td>
                  <td class="px-4 py-3">{{ formatTokenPrice(model.outputPrice) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section id="billing" class="doc-section">
          <h2>{{ docs.sections.billing }}</h2>
          <p>{{ docs.billingDescription }}</p>
        </section>

        <section id="errors" class="doc-section">
          <h2>{{ docs.sections.errors }}</h2>
          <ul>
            <li v-for="item in docs.errors" :key="item.code">
              <code>{{ item.code }}</code>: {{ item.description }}
            </li>
          </ul>
        </section>
      </article>

      <aside class="hidden lg:block">
        <div class="sticky top-24 rounded-lg border border-slate-200 bg-white p-4 text-sm dark:border-dark-800 dark:bg-dark-900">
          <h2 class="font-bold text-slate-900 dark:text-white">{{ docs.pageToc }}</h2>
          <div class="mt-3 space-y-2">
            <a
              v-for="item in docNav"
              :key="item.href"
              :href="item.href"
              class="block text-slate-500 hover:text-primary-700 dark:text-dark-300 dark:hover:text-primary-300"
            >
              {{ item.label }}
            </a>
          </div>
        </div>
      </aside>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import PublicTopNav from '@/components/public/PublicTopNav.vue'
import { dedupePublicModels, publicModels, normalizePublicModelLocale } from '@/constants/publicModels'
import { useAppStore } from '@/stores'

type DocsCopy = {
  kicker: string
  title: string
  description: string
  pageToc: string
  sections: Record<'quickstart' | 'auth' | 'models' | 'billing' | 'errors', string>
  quickstartSteps: string[]
  baseUrl: {
    current: string
    sdk: string
    httpsNote: string
  }
  authDescription: string
  chatDescription: string
  codexDescription: string
  codexConfigIntro: string
  authJsonExample: string
  modelsDescription: string
  billingDescription: string
  table: {
    model: string
    provider: string
    input: string
    output: string
  }
  errors: Array<{ code: string; description: string }>
}

type CopyKey = 'zh-CN' | 'zh-TW' | 'en' | 'th'

const DOCS_COPY: Record<CopyKey, DocsCopy> = {
  'zh-CN': {
    kicker: '开发文档',
    title: 'Token API Fuel 开发文档',
    description: '使用一个 API Key 调用平台已启用的主流 AI 模型。普通应用推荐使用 OpenAI Chat Completions 格式；Codex CLI 和 OpenCode 可使用本站兼容的 Responses 协议。',
    pageToc: '本页目录',
    sections: {
      quickstart: '快速开始',
      auth: '认证方式',
      models: '模型列表',
      billing: '计费与用量',
      errors: '常见错误',
    },
    quickstartSteps: [
      '注册并登录控制台。',
      '进入“API 密钥”创建一枚可用密钥。',
      '确认充值余额充足。',
      '按下方示例调用 /v1/chat/completions。',
    ],
    baseUrl: {
      current: '当前正式用户侧 API 地址：',
      sdk: 'OpenAI SDK 通常需要带 /v1 后缀：',
      httpsNote: '正式交付请优先使用 HTTPS。若临时排查证书或 443 端口问题，可短时间使用 HTTP。',
    },
    authDescription: '所有 API 请求都需要在 Header 中携带 Bearer Token：',
    chatDescription: '普通业务集成推荐使用该接口：',
    codexDescription: '本站兼容 /v1/responses、/responses 和 /backend-api/codex/responses 路由，因此 Codex CLI 示例保留 wire_api = "responses"。',
    codexConfigIntro: '请将以下内容放在 ~/.codex/config.toml 或 %userprofile%\\.codex\\config.toml 的开头：',
    authJsonExample: 'auth.json 示例：',
    modelsDescription: '推荐使用模型广场展示的模型名作为请求参数。当前常用模型包括：',
    billingDescription: '平台按 API Key 归属账户的充值余额扣减。每次成功调用会根据后台模型价格和 token usage 计算费用，并记录扣费金额和调用状态，可在控制台“使用记录”查看。',
    table: {
      model: '模型',
      provider: '供应商',
      input: '输入价格',
      output: '输出价格',
    },
    errors: [
      { code: '401', description: 'API Key 缺失、错误或已停用。' },
      { code: '402', description: '账户充值余额不足。' },
      { code: '404', description: '模型名称不存在或未启用。' },
      { code: '429', description: '请求过于频繁，请稍后重试。' },
      { code: '500', description: '网关或上游异常，请联系平台管理员。' },
    ],
  },
  'zh-TW': {
    kicker: '開發文件',
    title: 'Token API Fuel 開發文件',
    description: '使用一個 API Key 調用平台已啟用的主流 AI 模型。一般應用建議使用 OpenAI Chat Completions 格式；Codex CLI 和 OpenCode 可使用本站相容的 Responses 協議。',
    pageToc: '本頁目錄',
    sections: {
      quickstart: '快速開始',
      auth: '認證方式',
      models: '模型列表',
      billing: '計費與用量',
      errors: '常見錯誤',
    },
    quickstartSteps: [
      '註冊並登入控制台。',
      '進入「API 金鑰」建立一枚可用金鑰。',
      '確認儲值餘額充足。',
      '按下方範例調用 /v1/chat/completions。',
    ],
    baseUrl: {
      current: '目前正式使用者側 API 地址：',
      sdk: 'OpenAI SDK 通常需要帶 /v1 後綴：',
      httpsNote: '正式交付請優先使用 HTTPS。若臨時排查憑證或 443 連接埠問題，可短時間使用 HTTP。',
    },
    authDescription: '所有 API 請求都需要在 Header 中攜帶 Bearer Token：',
    chatDescription: '一般業務整合建議使用此接口：',
    codexDescription: '本站相容 /v1/responses、/responses 和 /backend-api/codex/responses 路由，因此 Codex CLI 範例保留 wire_api = "responses"。',
    codexConfigIntro: '請將以下內容放在 ~/.codex/config.toml 或 %userprofile%\\.codex\\config.toml 的開頭：',
    authJsonExample: 'auth.json 範例：',
    modelsDescription: '建議使用模型廣場展示的模型名稱作為請求參數。目前常用模型包括：',
    billingDescription: '平台按 API Key 歸屬帳戶的儲值餘額扣減。每次成功調用會根據後台模型價格和 token usage 計算費用，並記錄扣費金額和調用狀態，可在控制台「使用記錄」查看。',
    table: {
      model: '模型',
      provider: '供應商',
      input: '輸入價格',
      output: '輸出價格',
    },
    errors: [
      { code: '401', description: 'API Key 缺失、錯誤或已停用。' },
      { code: '402', description: '帳戶儲值餘額不足。' },
      { code: '404', description: '模型名稱不存在或未啟用。' },
      { code: '429', description: '請求過於頻繁，請稍後重試。' },
      { code: '500', description: '網關或上游異常，請聯絡平台管理員。' },
    ],
  },
  en: {
    kicker: 'Developer Documentation',
    title: 'Token API Fuel Developer Docs',
    description: 'Use one API key to call the enabled mainstream AI models on this gateway. Standard applications should use OpenAI Chat Completions; Codex CLI and OpenCode can use the compatible Responses protocol.',
    pageToc: 'On this page',
    sections: {
      quickstart: 'Quick Start',
      auth: 'Authentication',
      models: 'Model List',
      billing: 'Billing and Usage',
      errors: 'Common Errors',
    },
    quickstartSteps: [
      'Create an account and sign in to the console.',
      'Open API Keys and create an active key.',
      'Make sure your top-up balance is sufficient.',
      'Call /v1/chat/completions using the examples below.',
    ],
    baseUrl: {
      current: 'Current production API base URL for users:',
      sdk: 'OpenAI SDKs usually require the /v1 suffix:',
      httpsNote: 'Use HTTPS for production delivery. HTTP should only be used temporarily while diagnosing certificate or port 443 issues.',
    },
    authDescription: 'Every API request must include a Bearer token in the request headers:',
    chatDescription: 'Use this endpoint for normal business integrations:',
    codexDescription: 'This gateway is compatible with /v1/responses, /responses, and /backend-api/codex/responses, so the Codex CLI example keeps wire_api = "responses".',
    codexConfigIntro: 'Place the following at the beginning of ~/.codex/config.toml or %userprofile%\\.codex\\config.toml:',
    authJsonExample: 'auth.json example:',
    modelsDescription: 'Use the model names shown in the model marketplace as request parameters. Common enabled models include:',
    billingDescription: 'Charges are deducted from the top-up balance of the account that owns the API key. Each successful request calculates cost from backend model pricing and token usage, then records the charged amount and request status. You can review details in Console > Usage.',
    table: {
      model: 'Model',
      provider: 'Provider',
      input: 'Input Price',
      output: 'Output Price',
    },
    errors: [
      { code: '401', description: 'The API key is missing, invalid, or disabled.' },
      { code: '402', description: 'The account top-up balance is insufficient.' },
      { code: '404', description: 'The model name does not exist or is not enabled.' },
      { code: '429', description: 'Too many requests. Please retry later.' },
      { code: '500', description: 'Gateway or upstream error. Contact the platform administrator.' },
    ],
  },
  th: {
    kicker: 'เอกสารสำหรับนักพัฒนา',
    title: 'เอกสารนักพัฒนา Token API Fuel',
    description: 'ใช้ API Key เดียวเพื่อเรียกใช้โมเดล AI ที่เปิดใช้งานบนเกตเวย์นี้ แอปทั่วไปแนะนำให้ใช้รูปแบบ OpenAI Chat Completions ส่วน Codex CLI และ OpenCode ใช้โปรโตคอล Responses ที่รองรับได้',
    pageToc: 'หัวข้อในหน้านี้',
    sections: {
      quickstart: 'เริ่มต้นอย่างรวดเร็ว',
      auth: 'การยืนยันตัวตน',
      models: 'รายการโมเดล',
      billing: 'การคิดเงินและการใช้งาน',
      errors: 'ข้อผิดพลาดที่พบบ่อย',
    },
    quickstartSteps: [
      'สมัครบัญชีและเข้าสู่ระบบคอนโซล',
      'ไปที่ API Keys แล้วสร้างคีย์ที่เปิดใช้งาน',
      'ตรวจสอบว่ายอดเงินเติมเงินเพียงพอ',
      'เรียก /v1/chat/completions ตามตัวอย่างด้านล่าง',
    ],
    baseUrl: {
      current: 'Base URL สำหรับ API ฝั่งผู้ใช้:',
      sdk: 'OpenAI SDK มักต้องใส่ /v1 ต่อท้าย:',
      httpsNote: 'สำหรับการใช้งานจริงควรใช้ HTTPS ส่วน HTTP ใช้ชั่วคราวเฉพาะตอนตรวจสอบปัญหาใบรับรองหรือพอร์ต 443',
    },
    authDescription: 'ทุกคำขอ API ต้องใส่ Bearer token ใน header:',
    chatDescription: 'ใช้ปลายทางนี้สำหรับการเชื่อมต่อธุรกิจทั่วไป:',
    codexDescription: 'เกตเวย์นี้รองรับ /v1/responses, /responses และ /backend-api/codex/responses ดังนั้นตัวอย่าง Codex CLI จึงใช้ wire_api = "responses"',
    codexConfigIntro: 'ใส่ค่าต่อไปนี้ไว้ที่ส่วนต้นของ ~/.codex/config.toml หรือ %userprofile%\\.codex\\config.toml:',
    authJsonExample: 'ตัวอย่าง auth.json:',
    modelsDescription: 'ใช้ชื่อโมเดลจาก Model Marketplace เป็นพารามิเตอร์ในคำขอ โมเดลที่เปิดใช้งานทั่วไปมีดังนี้:',
    billingDescription: 'ค่าใช้จ่ายจะหักจากยอดเติมเงินของบัญชีที่เป็นเจ้าของ API Key ทุกคำขอที่สำเร็จจะคำนวณจากราคาหลังบ้านของโมเดลและ token usage แล้วบันทึกยอดที่หักและสถานะคำขอ ตรวจสอบได้ที่ Console > Usage',
    table: {
      model: 'โมเดล',
      provider: 'ผู้ให้บริการ',
      input: 'ราคาอินพุต',
      output: 'ราคาเอาต์พุต',
    },
    errors: [
      { code: '401', description: 'API Key หายไป ไม่ถูกต้อง หรือถูกปิดใช้งาน' },
      { code: '402', description: 'ยอดเงินเติมเงินของบัญชีไม่เพียงพอ' },
      { code: '404', description: 'ไม่มีชื่อโมเดลนี้หรือยังไม่ได้เปิดใช้งาน' },
      { code: '429', description: 'ส่งคำขอบ่อยเกินไป โปรดลองใหม่ภายหลัง' },
      { code: '500', description: 'เกิดข้อผิดพลาดที่เกตเวย์หรือต้นทาง โปรดติดต่อผู้ดูแลแพลตฟอร์ม' },
    ],
  },
}

const models = dedupePublicModels(publicModels)
const appStore = useAppStore()
const { locale } = useI18n()

const SECTION_LABELS: Record<CopyKey, { baseUrl: string; chat: string; codex: string }> = {
  'zh-CN': {
    baseUrl: 'Base URL（接口地址）',
    chat: '对话补全接口',
    codex: 'Codex CLI / OpenCode 配置',
  },
  'zh-TW': {
    baseUrl: 'Base URL（接口地址）',
    chat: '對話補全接口',
    codex: 'Codex CLI / OpenCode 設定',
  },
  en: {
    baseUrl: 'Base URL',
    chat: 'Chat Completions',
    codex: 'Codex CLI / OpenCode Setup',
  },
  th: {
    baseUrl: 'Base URL ของ API',
    chat: 'Chat Completions',
    codex: 'ตั้งค่า Codex CLI / OpenCode',
  },
}

const copyKey = computed<CopyKey>(() => normalizePublicModelLocale(String(locale.value || 'zh-CN')))
const docs = computed(() => DOCS_COPY[copyKey.value])
const sectionLabels = computed(() => SECTION_LABELS[copyKey.value])

const docNav = computed(() => [
  { href: '#quickstart', label: docs.value.sections.quickstart },
  { href: '#base-url', label: sectionLabels.value.baseUrl },
  { href: '#auth', label: docs.value.sections.auth },
  { href: '#chat', label: sectionLabels.value.chat },
  { href: '#codex', label: sectionLabels.value.codex },
  { href: '#models', label: docs.value.sections.models },
  { href: '#billing', label: docs.value.sections.billing },
  { href: '#errors', label: docs.value.sections.errors },
])

const gatewayBaseUrl = computed(() => {
  const base = appStore.apiBaseUrl?.trim() || 'https://tokenapifuel.com'
  return base.replace(/^http:\/\/tokenapifuel\.com/i, 'https://tokenapifuel.com').replace(/\/+$/, '')
})

const apiBaseUrl = computed(() => `${gatewayBaseUrl.value}/v1`)

const examplePrompt = computed(() => {
  if (copyKey.value === 'th') return 'แนะนำตัวเองในหนึ่งประโยค'
  if (copyKey.value === 'en') return 'Introduce yourself in one sentence'
  if (copyKey.value === 'zh-TW') return '用一句話介紹你自己'
  return '用一句话介绍你自己'
})

const curlExample = computed(() => `curl ${apiBaseUrl.value}/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer sk-your-api-key" \\
  -d '{
    "model": "deepseek-v4-flash",
    "messages": [
      {"role": "user", "content": "${examplePrompt.value}"}
    ]
  }'`)

const pythonExample = computed(() => `from openai import OpenAI

client = OpenAI(
    api_key="sk-your-api-key",
    base_url="${apiBaseUrl.value}"
)

response = client.chat.completions.create(
    model="deepseek-v4-flash",
    messages=[{"role": "user", "content": "${examplePrompt.value}"}]
)

print(response.choices[0].message.content)`)

const nodeExample = computed(() => `import OpenAI from "openai";

const client = new OpenAI({
  apiKey: "sk-your-api-key",
  baseURL: "${apiBaseUrl.value}",
});

const response = await client.chat.completions.create({
  model: "deepseek-v4-flash",
  messages: [{ role: "user", content: "${examplePrompt.value}" }],
});

console.log(response.choices[0].message.content);`)

const goExample = computed(() => `config := openai.DefaultConfig("sk-your-api-key")
config.BaseURL = "${apiBaseUrl.value}"
client := openai.NewClientWithConfig(config)`)

const codexConfigExample = computed(() => `model_provider = "OpenAI"
model = "deepseek-v4-flash"
review_model = "deepseek-v4-flash"
disable_response_storage = true
network_access = "enabled"
windows_wsl_setup_acknowledged = true

[model_providers.OpenAI]
name = "OpenAI"
base_url = "${gatewayBaseUrl.value}"
wire_api = "responses"
requires_openai_auth = true`)

function formatTokenPrice(value?: number | null): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) return '-'
  return `฿${(value * 1_000_000).toFixed(4)} / 1M Tokens`
}
</script>

<style scoped>
.doc-section {
  @apply mt-8 rounded-lg border border-slate-200 bg-white p-6 shadow-sm dark:border-dark-800 dark:bg-dark-900;
}

.doc-section h2 {
  @apply text-2xl font-bold text-slate-950 dark:text-white;
}

.doc-section h3 {
  @apply mt-6 text-lg font-bold text-slate-900 dark:text-white;
}

.doc-section p,
.doc-section li {
  @apply mt-3 leading-7 text-slate-600 dark:text-dark-300;
}

.doc-section ol,
.doc-section ul {
  @apply mt-4 list-inside space-y-2;
}

.doc-section ol {
  @apply list-decimal;
}

.doc-section ul {
  @apply list-disc;
}

.doc-section pre {
  @apply mt-4 overflow-x-auto rounded-md bg-slate-950 p-4 text-sm leading-7 text-slate-50;
}

.doc-section code {
  @apply rounded bg-slate-100 px-1.5 py-0.5 text-sm text-slate-900 dark:bg-dark-800 dark:text-dark-100;
}

.doc-section pre code {
  @apply bg-transparent p-0 text-slate-50;
}
</style>
