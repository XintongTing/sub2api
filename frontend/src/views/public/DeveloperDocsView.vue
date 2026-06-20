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
          <p class="text-sm font-semibold text-primary-700 dark:text-primary-300">Developer Documentation</p>
          <h1 class="mt-2 text-3xl font-bold text-slate-950 dark:text-white">Token API Fuel 开发文档</h1>
          <p class="mt-4 max-w-3xl text-base leading-7 text-slate-600 dark:text-dark-200">
            使用一个 API Key 调用平台已启用的主流 AI 模型。普通应用推荐使用 OpenAI Chat Completions 格式；Codex CLI 和 OpenCode 可使用本站兼容的 Responses 协议。
          </p>
        </section>

        <section id="quickstart" class="doc-section">
          <h2>快速开始</h2>
          <ol>
            <li>注册并登录控制台。</li>
            <li>进入“API 密钥”创建一枚可用密钥。</li>
            <li>确认充值余额充足。</li>
            <li>按下方示例调用 <code>/v1/chat/completions</code>。</li>
          </ol>
        </section>

        <section id="base-url" class="doc-section">
          <h2>Base URL</h2>
          <p>当前正式用户侧 API 地址：</p>
          <pre><code>https://tokenapifuel.com</code></pre>
          <p>OpenAI SDK 通常需要带 <code>/v1</code> 后缀：</p>
          <pre><code>https://tokenapifuel.com/v1</code></pre>
          <p>如果临时排查 HTTPS，可以短时间使用 HTTP；正式交付请优先使用 HTTPS。</p>
        </section>

        <section id="auth" class="doc-section">
          <h2>认证方式</h2>
          <p>所有 API 请求都需要在 Header 中携带 Bearer Token：</p>
          <pre><code>Authorization: Bearer sk-your-api-key
Content-Type: application/json</code></pre>
        </section>

        <section id="chat" class="doc-section">
          <h2>Chat Completions</h2>
          <p>普通业务集成推荐使用该接口：</p>
          <pre><code>POST /v1/chat/completions</code></pre>

          <h3>cURL</h3>
          <pre><code>curl https://tokenapifuel.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "deepseek-v3.2",
    "messages": [
      {"role": "user", "content": "用一句话介绍你自己"}
    ]
  }'</code></pre>

          <h3>Python</h3>
          <pre><code>from openai import OpenAI

client = OpenAI(
    api_key="sk-your-api-key",
    base_url="https://tokenapifuel.com/v1"
)

response = client.chat.completions.create(
    model="deepseek-v3.2",
    messages=[{"role": "user", "content": "用一句话介绍你自己"}]
)

print(response.choices[0].message.content)</code></pre>

          <h3>Node.js</h3>
          <pre><code>import OpenAI from "openai";

const client = new OpenAI({
  apiKey: "sk-your-api-key",
  baseURL: "https://tokenapifuel.com/v1",
});

const response = await client.chat.completions.create({
  model: "deepseek-v3.2",
  messages: [{ role: "user", content: "用一句话介绍你自己" }],
});

console.log(response.choices[0].message.content);</code></pre>

          <h3>Go</h3>
          <pre><code>config := openai.DefaultConfig("sk-your-api-key")
config.BaseURL = "https://tokenapifuel.com/v1"
client := openai.NewClientWithConfig(config)</code></pre>
        </section>

        <section id="codex" class="doc-section">
          <h2>Codex CLI / OpenCode</h2>
          <p>本站已注册并兼容 <code>/v1/responses</code>、<code>/responses</code> 和 <code>/backend-api/codex/responses</code> 路由，因此 Codex CLI 示例保留 <code>wire_api = "responses"</code>。</p>
          <p>请将以下内容放在 <code>~/.codex/config.toml</code> 或 <code>%userprofile%\.codex\config.toml</code> 的开头：</p>
          <pre><code>model_provider = "OpenAI"
model = "deepseek-v3.2"
review_model = "deepseek-v3.2"
disable_response_storage = true
network_access = "enabled"
windows_wsl_setup_acknowledged = true

[model_providers.OpenAI]
name = "OpenAI"
base_url = "https://tokenapifuel.com"
wire_api = "responses"
requires_openai_auth = true</code></pre>

          <p><code>auth.json</code> 示例：</p>
          <pre><code>{
  "OPENAI_API_KEY": "sk-your-api-key"
}</code></pre>
        </section>

        <section id="models" class="doc-section">
          <h2>模型列表</h2>
          <p>推荐使用模型广场展示的模型名作为请求参数。当前常用模型包括：</p>
          <div class="overflow-hidden rounded-lg border border-slate-200 dark:border-dark-800">
            <table class="min-w-full divide-y divide-slate-200 text-sm dark:divide-dark-800">
              <thead class="bg-slate-50 text-left font-semibold text-slate-600 dark:bg-dark-800 dark:text-dark-200">
                <tr>
                  <th class="px-4 py-3">模型</th>
                  <th class="px-4 py-3">供应商</th>
                  <th class="px-4 py-3">输入价格</th>
                  <th class="px-4 py-3">输出价格</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-dark-800">
                <tr v-for="model in models" :key="model.id">
                  <td class="px-4 py-3 font-semibold">{{ model.displayName }}</td>
                  <td class="px-4 py-3">{{ model.provider }}</td>
                  <td class="px-4 py-3">{{ model.inputPrice }} / {{ model.unit }}</td>
                  <td class="px-4 py-3">{{ model.outputPrice }} / {{ model.unit }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section id="billing" class="doc-section">
          <h2>计费与用量</h2>
          <p>平台按用户 API Key 归属账户的充值余额扣减。每次成功调用会根据模型价格和 Token usage 计算费用，并记录扣费金额和调用状态，可在控制台“使用记录”查看。</p>
        </section>

        <section id="errors" class="doc-section">
          <h2>常见错误</h2>
          <ul>
            <li><code>401</code>：API Key 缺失、错误或已停用。</li>
            <li><code>402</code>：账户充值余额不足。</li>
            <li><code>404</code>：模型名称不存在或未启用。</li>
            <li><code>429</code>：请求过于频繁，请稍后重试。</li>
            <li><code>500</code>：网关或上游异常，请联系平台管理员。</li>
          </ul>
        </section>
      </article>

      <aside class="hidden lg:block">
        <div class="sticky top-24 rounded-lg border border-slate-200 bg-white p-4 text-sm dark:border-dark-800 dark:bg-dark-900">
          <h2 class="font-bold text-slate-900 dark:text-white">本页目录</h2>
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
import PublicTopNav from '@/components/public/PublicTopNav.vue'
import { publicModels } from '@/constants/publicModels'

const models = publicModels
const docNav = [
  { href: '#quickstart', label: '快速开始' },
  { href: '#base-url', label: 'Base URL' },
  { href: '#auth', label: '认证方式' },
  { href: '#chat', label: 'Chat Completions' },
  { href: '#codex', label: 'Codex / OpenCode' },
  { href: '#models', label: '模型列表' },
  { href: '#billing', label: '计费与用量' },
  { href: '#errors', label: '常见错误' },
]
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
