<template>
  <div class="min-h-screen bg-slate-50 text-slate-950 dark:bg-dark-950 dark:text-white">
    <PublicTopNav />

    <main class="mx-auto max-w-6xl px-4 py-10 sm:px-6 lg:px-8 lg:py-14">
      <section class="overflow-hidden rounded-2xl bg-slate-950 text-white shadow-xl dark:ring-1 dark:ring-white/10">
        <div class="grid gap-8 px-6 py-10 sm:px-10 lg:grid-cols-[1.25fr_0.75fr] lg:px-12 lg:py-12">
          <div>
            <p class="text-sm font-semibold uppercase tracking-[0.18em] text-cyan-300">{{ copy.kicker }}</p>
            <h1 class="mt-4 max-w-3xl text-3xl font-bold leading-tight sm:text-5xl">{{ copy.title }}</h1>
            <p class="mt-5 max-w-2xl text-base leading-7 text-slate-300">{{ copy.description }}</p>
          </div>
          <div class="rounded-xl border border-white/15 bg-white/10 p-5 backdrop-blur">
            <div class="flex items-start gap-3">
              <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-cyan-400/15 text-cyan-200">
                <Icon name="mail" size="md" />
              </span>
              <div>
                <p class="text-sm text-slate-300">{{ copy.humanSupport }}</p>
                <a :href="generalMailto" class="mt-1 block break-all font-semibold text-white underline decoration-cyan-300 underline-offset-4">
                  {{ supportEmail }}
                </a>
              </div>
            </div>
            <p class="mt-4 text-sm leading-6 text-slate-300">{{ copy.responseNote }}</p>
          </div>
        </div>
      </section>

      <section class="mt-8 grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
        <div class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-dark-800 dark:bg-dark-900 sm:p-8">
          <p class="text-sm font-semibold text-primary-700 dark:text-primary-300">{{ copy.autoReply }}</p>
          <h2 class="mt-2 text-2xl font-bold">{{ copy.searchTitle }}</h2>
          <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-dark-300">{{ copy.searchDescription }}</p>

          <label class="mt-6 block">
            <span class="sr-only">{{ copy.searchPlaceholder }}</span>
            <div class="relative">
              <Icon name="search" size="sm" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
              <input
                v-model="query"
                type="search"
                :placeholder="copy.searchPlaceholder"
                class="w-full rounded-xl border border-slate-300 bg-slate-50 py-3 pl-11 pr-4 text-sm outline-none transition focus:border-primary-500 focus:bg-white focus:ring-4 focus:ring-primary-500/10 dark:border-dark-700 dark:bg-dark-950 dark:focus:border-primary-400"
              />
            </div>
          </label>

          <div class="mt-4 flex flex-wrap gap-2">
            <button
              v-for="keyword in copy.keywords"
              :key="keyword"
              type="button"
              class="rounded-full border border-slate-200 bg-white px-3 py-1.5 text-sm font-semibold text-slate-600 transition hover:border-primary-300 hover:bg-primary-50 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:bg-primary-500/10 dark:hover:text-primary-300"
              @click="query = keyword"
            >
              {{ keyword }}
            </button>
          </div>

          <div class="mt-6 rounded-xl border border-primary-100 bg-primary-50/70 p-5 dark:border-primary-500/20 dark:bg-primary-500/10">
            <div class="flex items-start gap-3">
              <Icon name="chat" size="md" class="mt-0.5 shrink-0 text-primary-700 dark:text-primary-300" />
              <div>
                <h3 class="font-bold text-slate-950 dark:text-white">{{ activeAnswer.title }}</h3>
                <p class="mt-2 whitespace-pre-line text-sm leading-6 text-slate-700 dark:text-dark-200">{{ activeAnswer.answer }}</p>
                <router-link
                  v-if="activeAnswer.path"
                  :to="activeAnswer.path"
                  class="mt-4 inline-flex items-center text-sm font-bold text-primary-700 hover:text-primary-800 dark:text-primary-300"
                >
                  {{ activeAnswer.linkLabel }} →
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <div class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-dark-800 dark:bg-dark-900 sm:p-8">
          <p class="text-sm font-semibold text-primary-700 dark:text-primary-300">{{ copy.emailKicker }}</p>
          <h2 class="mt-2 text-2xl font-bold">{{ copy.emailTitle }}</h2>
          <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-dark-300">{{ copy.emailDescription }}</p>

          <div class="mt-6 space-y-3">
            <a
              v-for="item in mailTopics"
              :key="item.subject"
              :href="item.href"
              class="group flex items-center justify-between gap-4 rounded-xl border border-slate-200 p-4 transition hover:border-primary-300 hover:bg-primary-50/60 dark:border-dark-700 dark:hover:bg-primary-500/10"
            >
              <span>
                <span class="block font-bold text-slate-900 dark:text-white">{{ item.label }}</span>
                <span class="mt-1 block text-sm text-slate-500 dark:text-dark-300">{{ item.description }}</span>
              </span>
              <Icon name="arrowRight" size="sm" class="shrink-0 text-slate-400 transition group-hover:translate-x-1 group-hover:text-primary-600" />
            </a>
          </div>

          <div class="mt-6 rounded-xl bg-slate-100 p-4 text-sm leading-6 text-slate-600 dark:bg-dark-800 dark:text-dark-200">
            <strong class="text-slate-900 dark:text-white">{{ copy.securityTitle }}</strong>
            {{ copy.securityNote }}
          </div>
        </div>
      </section>

      <div class="mt-8 text-center text-xs text-slate-400 dark:text-dark-500">
        <a href="https://deerflow.tech" target="_blank" rel="noopener noreferrer" class="transition hover:text-slate-600 dark:hover:text-dark-300">
          Created By Deerflow
        </a>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PublicTopNav from '@/components/public/PublicTopNav.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores'

type Answer = { title: string; answer: string; path?: string; linkLabel?: string }

const DEFAULT_SUPPORT_EMAIL = 'service@tokenapifuel.com'
const { locale } = useI18n()
const appStore = useAppStore()
const query = ref('')

const supportEmail = computed(() => {
  const configured = String(appStore.contactInfo || '').trim()
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(configured) ? configured : DEFAULT_SUPPORT_EMAIL
})

const COPY = {
  zh: {
    kicker: '客服中心', title: '先自助定位，再由人工继续处理',
    description: '输入问题关键词可立即获得操作指引；账户、支付、退款及隐私事项可通过注册邮箱提交给人工客服。',
    humanSupport: '人工客服邮箱', responseNote: '建议使用账户注册邮箱联系，并在邮件中提供订单号或问题发生时间。',
    autoReply: '关键词自动回复', searchTitle: '你遇到了什么问题？', searchDescription: '支持“充值、支付失败、未到账、退款、API Key、余额、发票、隐私”等关键词。',
    searchPlaceholder: '例如：支付成功但余额未到账', keywords: ['充值', '支付失败', '未到账', '退款', 'API Key', '余额', '隐私'],
    emailKicker: '人工客服', emailTitle: '按事项发送邮件', emailDescription: '选择事项后会自动填写邮件主题，便于客服分类处理。',
    emailTopics: [
      ['支付或入账问题', '请附订单号、时间、金额、币种及脱敏凭证', 'Payment Support / 支付问题'],
      ['退款申请', '请说明退款原因并附订单信息', 'Refund Request / 退款申请'],
      ['账户或 API Key', '请勿在邮件中发送密码或完整 API Key', 'Account and API Key Support / 账户与API Key'],
      ['隐私请求', '访问、更正、删除或其他个人资料请求', 'Privacy Request / 隐私请求'],
    ],
    securityTitle: '安全提示：', securityNote: '客服不会索要密码、完整银行卡号、银行卡安全码或秘密 API Key。',
    defaultAnswer: { title: '请输入问题关键词', answer: '我们会根据关键词显示对应操作。若未找到答案，请使用右侧人工客服入口发送邮件。' },
    answers: {
      payment: { title: '充值与支付', answer: '登录后进入充值页面，选择金额和当前可用的支付方式，核对币种与最终金额后完成支付。支付状态可在支付结果页或订单记录中查看。', path: '/legal/payment-process', linkLabel: '查看完整支付流程' },
      pending: { title: '支付成功但余额未到账', answer: '先刷新支付结果页并查看订单记录。若款项已扣除且仍未入账，请使用注册邮箱联系客服，提供订单号、支付时间、金额、币种及已遮盖敏感信息的支付凭证。', path: '/legal/delivery-policy', linkLabel: '查看数字交付政策' },
      refund: { title: '退款申请', answer: '使用注册邮箱发送邮件，主题注明“Refund Request / 退款申请”，提供订单号、日期、金额、币种、申请原因及脱敏证明。已消耗的数字服务余额通常不支持退款，法定权利除外。', path: '/legal/refund-policy', linkLabel: '查看退款政策' },
      api: { title: 'API Key 与调用问题', answer: '登录控制台创建或管理 API Key。遇到 401 请检查 Key 是否正确或停用；402 通常表示余额不足；404 请核对模型名称；429 表示请求过于频繁。', path: '/docs', linkLabel: '查看开发文档' },
      privacy: { title: '隐私与个人资料请求', answer: '请使用注册邮箱联系客服，主题注明“Privacy Request / 隐私请求”，说明希望访问、更正、删除或限制处理的资料。客服可能需要核实身份。', path: '/legal/privacy', linkLabel: '查看隐私政策' },
    } as Record<string, Answer>,
  },
  en: {
    kicker: 'Customer support', title: 'Find an answer first, then continue with a person',
    description: 'Enter a keyword for immediate guidance. Account, payment, refund, and privacy matters can be sent to human support from your registered email address.',
    humanSupport: 'Human support email', responseNote: 'Use your registered email address and include the order number or time the issue occurred.',
    autoReply: 'Keyword auto-reply', searchTitle: 'What can we help with?', searchDescription: 'Try top up, payment failed, not credited, refund, API key, balance, invoice, or privacy.',
    searchPlaceholder: 'Example: payment completed but balance not credited', keywords: ['top up', 'payment failed', 'not credited', 'refund', 'API key', 'balance', 'privacy'],
    emailKicker: 'Human support', emailTitle: 'Email by topic', emailDescription: 'Choose a topic to prefill the subject for faster classification.',
    emailTopics: [
      ['Payment or delivery', 'Include order, time, amount, currency, and a redacted receipt', 'Payment Support'],
      ['Refund request', 'Include the order details and reason for the request', 'Refund Request'],
      ['Account or API key', 'Never include a password or full API key', 'Account and API Key Support'],
      ['Privacy request', 'Access, correction, deletion, or other data requests', 'Privacy Request'],
    ],
    securityTitle: 'Security:', securityNote: 'Support will never ask for your password, full card number, card security code, or secret API key.',
    defaultAnswer: { title: 'Enter a question keyword', answer: 'We will show the relevant steps. If no answer matches, use the human-support email options.' },
    answers: {
      payment: { title: 'Top-up and payment', answer: 'Sign in, open the purchase page, select an amount and an available payment method, then review the currency and final amount before confirming. Check the payment-result page or order history for status.', path: '/legal/payment-process', linkLabel: 'View payment process' },
      pending: { title: 'Paid but not credited', answer: 'Refresh the payment-result page and check order history. If the charge is confirmed but no balance appears, email support from the registered address with the order number, time, amount, currency, and a redacted receipt.', path: '/legal/delivery-policy', linkLabel: 'View delivery policy' },
      refund: { title: 'Refund request', answer: 'Email from the registered address with subject “Refund Request”, order number, date, amount, currency, reason, and redacted evidence. Consumed digital-service balance is normally non-refundable except where required by law.', path: '/legal/refund-policy', linkLabel: 'View refund policy' },
      api: { title: 'API key and request issues', answer: 'Create or manage keys in the console. A 401 usually means a missing or invalid key; 402 means insufficient balance; 404 means the model is unavailable; 429 means too many requests.', path: '/docs', linkLabel: 'View developer docs' },
      privacy: { title: 'Privacy and data requests', answer: 'Email from the registered address with subject “Privacy Request” and describe whether you seek access, correction, deletion, or restriction. Identity verification may be required.', path: '/legal/privacy', linkLabel: 'View privacy policy' },
    } as Record<string, Answer>,
  },
}

const copy = computed(() => String(locale.value || '').toLowerCase().startsWith('zh') ? COPY.zh : COPY.en)
const normalizedQuery = computed(() => query.value.trim().toLowerCase())

const activeAnswer = computed<Answer>(() => {
  const value = normalizedQuery.value
  if (!value) return copy.value.defaultAnswer
  if (/退款|refund|退费/.test(value)) return copy.value.answers.refund
  if (/未到账|没到账|not credited|pending|扣款|charged/.test(value)) return copy.value.answers.pending
  if (/api|key|密钥|401|402|404|429|余额|balance/.test(value)) return copy.value.answers.api
  if (/隐私|privacy|删除数据|个人资料/.test(value)) return copy.value.answers.privacy
  if (/充值|支付|付款|top.?up|payment|invoice|发票/.test(value)) return copy.value.answers.payment
  return copy.value.defaultAnswer
})

const generalMailto = computed(() => `mailto:${supportEmail.value}?subject=${encodeURIComponent('Customer Support / 客户服务')}`)
const mailTopics = computed(() => copy.value.emailTopics.map(([label, description, subject]) => ({
  label,
  description,
  subject,
  href: `mailto:${supportEmail.value}?subject=${encodeURIComponent(subject)}`,
})))
</script>
