<template>
  <div class="grid min-h-[calc(100vh-4rem)] bg-white text-slate-950 dark:bg-dark-950 dark:text-white xl:grid-cols-[340px_1fr]">
    <aside class="border-b border-slate-200 bg-white p-4 dark:border-dark-800 dark:bg-dark-950 xl:border-b-0 xl:border-r">
      <div class="mb-5 flex items-center gap-3">
        <span class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary-50 text-primary-700 dark:bg-primary-500/10 dark:text-primary-300">
          <Icon name="terminal" size="md" />
        </span>
        <div>
          <h1 class="text-lg font-bold">模型配置</h1>
          <p class="text-sm text-slate-500 dark:text-dark-300">API Key 只保存在当前页面状态</p>
        </div>
      </div>

      <form class="space-y-5" @submit.prevent="sendMessage">
        <div>
          <label for="api-key" class="config-label">API Key</label>
          <input
            id="api-key"
            v-model.trim="apiKey"
            type="password"
            autocomplete="off"
            placeholder="sk-..."
            class="config-input"
          />
        </div>

        <div>
          <label for="group" class="config-label">分组</label>
          <select id="group" v-model="group" class="config-input">
            <option value="default">默认分组</option>
            <option value="vip">vip</option>
          </select>
        </div>

        <div>
          <label for="model" class="config-label">模型</label>
          <select id="model" v-model="selectedModel" class="config-input">
            <option v-for="model in businessModels" :key="model.displayName" :value="model.displayName">
              {{ model.displayName }}
            </option>
          </select>
        </div>

        <div class="rounded-lg border border-slate-200 p-3 dark:border-dark-800">
          <label class="flex items-center justify-between gap-3 text-sm font-semibold">
            <span>自定义请求体模式</span>
            <input v-model="customRequestEnabled" type="checkbox" class="h-4 w-4 accent-primary-600" />
          </label>
          <textarea
            v-if="customRequestEnabled"
            v-model="customRequestJson"
            rows="7"
            class="mt-3 config-textarea font-mono text-xs"
            placeholder='{"max_tokens": 512}'
          ></textarea>
        </div>

        <div class="rounded-lg border border-slate-200 p-3 dark:border-dark-800">
          <label class="flex items-center justify-between gap-3 text-sm font-semibold">
            <span>图片地址</span>
            <input v-model="imageUrlEnabled" type="checkbox" class="h-4 w-4 accent-primary-600" />
          </label>
          <input
            v-if="imageUrlEnabled"
            v-model.trim="imageUrl"
            type="url"
            class="mt-3 config-input"
            placeholder="https://example.com/image.jpg"
          />
        </div>

        <ParameterSlider label="Temperature" help="控制输出的随机性和创造性" :min="0" :max="2" :step="0.1" v-model="temperature" />
        <ParameterSlider label="Top P" help="核采样，控制词汇选择的多样性" :min="0" :max="1" :step="0.05" v-model="topP" />
        <ParameterSlider label="Frequency Penalty" help="减少重复词汇的出现" :min="-2" :max="2" :step="0.1" v-model="frequencyPenalty" />
        <ParameterSlider label="Presence Penalty" help="鼓励讨论新的主题" :min="-2" :max="2" :step="0.1" v-model="presencePenalty" />
      </form>
    </aside>

    <section class="flex min-h-0 flex-col bg-slate-50/80 dark:bg-dark-950">
      <header class="flex items-center justify-between border-b border-slate-200 bg-white px-5 py-4 dark:border-dark-800 dark:bg-dark-950">
        <div class="flex min-w-0 items-center gap-3">
          <Icon name="chat" size="md" class="shrink-0 text-primary-600" />
          <div class="min-w-0">
            <h2 class="truncate text-lg font-bold">AI 对话</h2>
            <p class="truncate text-sm text-slate-500 dark:text-dark-300">{{ selectedModel }}</p>
          </div>
        </div>
        <button type="button" class="toolbar-button" @click="showRaw = !showRaw">
          <Icon name="eye" size="sm" />
          {{ showRaw ? '隐藏调试' : '显示调试' }}
        </button>
      </header>

      <div class="min-h-0 flex-1 overflow-y-auto p-5">
        <div class="mx-auto max-w-5xl space-y-5">
          <div class="flex gap-3">
            <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-primary-50 text-primary-700 dark:bg-primary-500/10 dark:text-primary-300">
              <Icon name="sparkles" size="sm" />
            </div>
            <div class="rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm leading-6 shadow-sm dark:border-dark-800 dark:bg-dark-900">
              你好！粘贴你的 API Key，选择模型，然后发送一条消息。我会通过本站网关发起真实调用，并展示 usage 和原始响应。
            </div>
          </div>

          <div v-for="message in messages" :key="message.id" class="flex gap-3" :class="{ 'justify-end': message.role === 'user' }">
            <div v-if="message.role === 'assistant'" class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-primary-50 text-primary-700 dark:bg-primary-500/10 dark:text-primary-300">
              <Icon name="sparkles" size="sm" />
            </div>
            <div
              class="max-w-[min(760px,85%)] rounded-lg px-4 py-3 text-sm leading-6 shadow-sm"
              :class="message.role === 'user'
                ? 'bg-primary-600 text-white'
                : 'border border-slate-200 bg-white text-slate-800 dark:border-dark-800 dark:bg-dark-900 dark:text-dark-100'"
            >
              <pre class="whitespace-pre-wrap font-sans">{{ message.content }}</pre>
            </div>
            <div v-if="message.role === 'user'" class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-primary-600 text-white">
              我
            </div>
          </div>

          <div v-if="errorMessage" class="rounded-lg border border-red-200 bg-red-50 p-4 text-sm leading-6 text-red-700 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-200">
            {{ errorMessage }}
          </div>

          <div v-if="usage.total_tokens !== undefined" class="grid gap-3 sm:grid-cols-3">
            <UsageCard label="Prompt Tokens" :value="usage.prompt_tokens" />
            <UsageCard label="Completion Tokens" :value="usage.completion_tokens" />
            <UsageCard label="Total Tokens" :value="usage.total_tokens" />
          </div>

          <details v-if="showRaw && rawResponse" open class="rounded-lg border border-slate-200 bg-white dark:border-dark-800 dark:bg-dark-900">
            <summary class="cursor-pointer px-4 py-3 text-sm font-semibold">原始响应 JSON</summary>
            <pre class="overflow-x-auto border-t border-slate-200 p-4 text-xs leading-5 dark:border-dark-800"><code>{{ rawResponse }}</code></pre>
          </details>
        </div>
      </div>

      <form class="border-t border-slate-200 bg-white p-4 dark:border-dark-800 dark:bg-dark-950" @submit.prevent="sendMessage">
        <div class="mx-auto flex max-w-5xl items-end gap-3 rounded-lg border border-slate-200 bg-white p-2 shadow-sm dark:border-dark-800 dark:bg-dark-900">
          <button type="button" class="rounded-md p-2 text-slate-400 hover:bg-slate-100 hover:text-slate-700 dark:hover:bg-dark-800" @click="clearChat" title="清空对话">
            <Icon name="trash" size="md" />
          </button>
          <textarea
            v-model="prompt"
            rows="1"
            class="min-h-10 flex-1 resize-none bg-transparent px-2 py-2 text-sm outline-none"
            placeholder="请输入您的问题..."
            @keydown.enter.exact.prevent="sendMessage"
          ></textarea>
          <button
            type="submit"
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary-600 text-white shadow-sm shadow-primary-600/20 transition hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="isLoading || !canSubmit"
            title="发送"
          >
            <Icon v-if="isLoading" name="refresh" size="md" class="animate-spin" />
            <Icon v-else name="arrowUp" size="md" />
          </button>
        </div>
      </form>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, reactive, ref } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import { dedupePublicModels, publicModels } from '@/constants/publicModels'

interface ChatMessage {
  id: number
  role: 'user' | 'assistant'
  content: string
}

interface ChatCompletionResponse {
  choices?: Array<{
    message?: { content?: string }
    text?: string
  }>
  usage?: {
    prompt_tokens?: number
    completion_tokens?: number
    total_tokens?: number
  }
  error?: { message?: string }
}

const ParameterSlider = defineComponent({
  props: {
    label: { type: String, required: true },
    help: { type: String, required: true },
    min: { type: Number, required: true },
    max: { type: Number, required: true },
    step: { type: Number, required: true },
    modelValue: { type: Number, required: true },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    return () => h('div', { class: 'space-y-2' }, [
      h('div', { class: 'flex items-center justify-between gap-3' }, [
        h('label', { class: 'text-sm font-semibold text-slate-900 dark:text-white' }, props.label),
        h('span', { class: 'rounded-full bg-slate-100 px-2 py-0.5 text-xs font-semibold text-slate-600 dark:bg-dark-800 dark:text-dark-200' }, String(props.modelValue)),
      ]),
      h('input', {
        type: 'range',
        min: props.min,
        max: props.max,
        step: props.step,
        value: props.modelValue,
        class: 'w-full accent-primary-600',
        onInput: (event: Event) => emit('update:modelValue', Number((event.target as HTMLInputElement).value)),
      }),
      h('p', { class: 'text-xs text-slate-500 dark:text-dark-400' }, props.help),
    ])
  },
})

const UsageCard = defineComponent({
  props: {
    label: { type: String, required: true },
    value: { type: Number, default: undefined },
  },
  setup(props) {
    return () => h('div', { class: 'rounded-lg border border-slate-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900' }, [
      h('div', { class: 'text-xs font-semibold uppercase text-slate-500 dark:text-dark-300' }, props.label),
      h('div', { class: 'mt-2 text-2xl font-bold text-primary-700 dark:text-primary-300' }, props.value ?? '-'),
    ])
  },
})

const businessModels = dedupePublicModels(publicModels).filter(model => model.endpointTypes.includes('openai:/v1/chat/completions'))
const apiKey = ref('')
const group = ref('default')
const selectedModel = ref(businessModels.find(model => model.displayName === 'deepseek-v4-flash')?.displayName || businessModels[0]?.displayName || 'deepseek-v4-flash')
const customRequestEnabled = ref(false)
const customRequestJson = ref('{\n  "max_tokens": 512\n}')
const imageUrlEnabled = ref(false)
const imageUrl = ref('')
const temperature = ref(0.7)
const topP = ref(1)
const frequencyPenalty = ref(0)
const presencePenalty = ref(0)
const prompt = ref('')
const messages = ref<ChatMessage[]>([])
const errorMessage = ref('')
const rawResponse = ref('')
const showRaw = ref(false)
const isLoading = ref(false)
let nextMessageId = 1

const usage = reactive({
  prompt_tokens: undefined as number | undefined,
  completion_tokens: undefined as number | undefined,
  total_tokens: undefined as number | undefined,
})

const canSubmit = computed(() => apiKey.value.length > 8 && selectedModel.value && prompt.value.trim().length > 0)

function clearResult(): void {
  errorMessage.value = ''
  rawResponse.value = ''
  usage.prompt_tokens = undefined
  usage.completion_tokens = undefined
  usage.total_tokens = undefined
}

function clearChat(): void {
  messages.value = []
  clearResult()
}

function buildMessageContent(text: string): unknown {
  if (!imageUrlEnabled.value || !imageUrl.value) {
    return text
  }
  return [
    { type: 'text', text },
    { type: 'image_url', image_url: { url: imageUrl.value } },
  ]
}

async function sendMessage(): Promise<void> {
  if (!canSubmit.value || isLoading.value) return

  const text = prompt.value.trim()
  prompt.value = ''
  clearResult()
  messages.value.push({ id: nextMessageId++, role: 'user', content: text })
  isLoading.value = true

  try {
    const body: Record<string, unknown> = {
      model: selectedModel.value,
      messages: [{ role: 'user', content: buildMessageContent(text) }],
      temperature: temperature.value,
      top_p: topP.value,
      frequency_penalty: frequencyPenalty.value,
      presence_penalty: presencePenalty.value,
    }

    if (customRequestEnabled.value) {
      const extra = JSON.parse(customRequestJson.value || '{}') as Record<string, unknown>
      Object.assign(body, extra)
    }

    const response = await fetch('/v1/chat/completions', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${apiKey.value}`,
        'Content-Type': 'application/json',
        'X-TokenAPI-Group': group.value,
      },
      body: JSON.stringify(body),
    })

    const data = (await response.json().catch(() => ({}))) as ChatCompletionResponse
    rawResponse.value = JSON.stringify(data, null, 2)
    if (!response.ok) {
      throw new Error(data.error?.message || `请求失败，HTTP ${response.status}`)
    }

    const content = data.choices?.[0]?.message?.content || data.choices?.[0]?.text || '调用成功，但响应内容为空。'
    messages.value.push({ id: nextMessageId++, role: 'assistant', content })
    usage.prompt_tokens = data.usage?.prompt_tokens
    usage.completion_tokens = data.usage?.completion_tokens
    usage.total_tokens = data.usage?.total_tokens
  } catch (error) {
    const message = error instanceof SyntaxError
      ? '自定义请求体不是合法 JSON，请检查后重试。'
      : error instanceof Error
        ? error.message
        : '请求失败，请检查 API Key、余额、模型名称和服务状态。'
    errorMessage.value = message
    messages.value.push({ id: nextMessageId++, role: 'assistant', content: message })
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.config-label {
  @apply text-sm font-semibold text-slate-900 dark:text-white;
}

.config-input {
  @apply mt-2 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 dark:border-dark-700 dark:bg-dark-950 dark:text-white;
}

.config-textarea {
  @apply w-full resize-y rounded-md border border-slate-300 bg-white px-3 py-2 text-sm leading-6 text-slate-900 outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 dark:border-dark-700 dark:bg-dark-950 dark:text-white;
}

.toolbar-button {
  @apply inline-flex items-center gap-1.5 rounded-md border border-slate-200 bg-white px-3 py-2 text-sm font-semibold text-slate-600 transition hover:border-primary-200 hover:bg-primary-50 hover:text-primary-700 dark:border-dark-800 dark:bg-dark-900 dark:text-dark-200 dark:hover:bg-primary-500/10;
}
</style>
