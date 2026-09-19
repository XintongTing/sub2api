<template>
  <div>
    <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
      {{ t('payment.paymentMethod') }}
    </label>
    <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
      <button
        v-for="method in displayMethods"
        :key="`${method.type}-${method.label || method.type}`"
        type="button"
        :disabled="!method.available"
        :class="[
          'relative flex h-[60px] flex-col items-center justify-center rounded-lg border px-3 transition-all sm:flex-1',
          !method.available
            ? 'cursor-not-allowed border-gray-200 bg-gray-50 opacity-50 dark:border-dark-700 dark:bg-dark-800/50'
            : selected === method.type
              ? methodSelectedClass(method.type)
              : 'border-gray-300 bg-white text-gray-700 hover:border-gray-400 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-200 dark:hover:border-dark-500',
        ]"
        @click="method.available && emit('select', method.type)"
      >
        <span class="flex items-center gap-2">
          <img :src="method.icon || methodIcon(method.type)" :alt="method.label || t(`payment.methods.${method.type}`)" class="h-7 w-7 object-contain" />
          <span class="flex flex-col items-start leading-none">
            <span class="text-base font-semibold">{{ method.label || t(`payment.methods.${method.type}`) }}</span>
            <span
              v-if="method.fee_rate > 0"
              class="text-[10px] tracking-wide text-gray-500 dark:text-dark-400"
            >
              {{ t('payment.fee') }} {{ method.fee_rate }}%
            </span>
          </span>
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { METHOD_ORDER } from './providerConfig'
import alipayIcon from '@/assets/icons/alipay.svg'
import wxpayIcon from '@/assets/icons/wxpay.svg'
import stripeIcon from '@/assets/icons/stripe.svg'
import airwallexIcon from '@/assets/icons/airwallex.svg'
import promptpayIcon from '@/assets/icons/promptpay.jpg'
import kplusIcon from '@/assets/icons/kplus.jpg'
import truemoneyIcon from '@/assets/icons/truemoney.jpg'
import linepayIcon from '@/assets/icons/linepay.jpg'

export interface PaymentMethodOption {
  type: string
  fee_rate: number
  available: boolean
  label?: string
  icon?: string
}

const props = defineProps<{
  methods: PaymentMethodOption[]
  selected: string
}>()

const emit = defineEmits<{
  select: [type: string]
}>()

const { t } = useI18n()

const METHOD_ICONS: Record<string, string> = {
  alipay: alipayIcon,
  wxpay: wxpayIcon,
  stripe: stripeIcon,
  airwallex: airwallexIcon,
}

const SUNRATE_METHODS = [
  { type: 'sunrate', label: 'PromptPay', icon: promptpayIcon },
  { type: 'sunrate', label: 'K PLUS', icon: kplusIcon },
  { type: 'sunrate', label: 'TrueMoney', icon: truemoneyIcon },
  { type: 'sunrate', label: 'LINE Pay', icon: linepayIcon },
]

const sortedMethods = computed(() => {
  const order: readonly string[] = METHOD_ORDER
  return [...props.methods].sort((a, b) => {
    const ai = order.indexOf(a.type)
    const bi = order.indexOf(b.type)
    return (ai === -1 ? 999 : ai) - (bi === -1 ? 999 : bi)
  })
})

const displayMethods = computed(() => sortedMethods.value.flatMap(method =>
  method.type === 'sunrate'
    ? SUNRATE_METHODS.map(local => ({ ...method, ...local }))
    : [method],
))

function methodIcon(type: string): string {
  if (type.includes('alipay')) return METHOD_ICONS.alipay
  if (type.includes('wxpay')) return METHOD_ICONS.wxpay
  if (type === 'airwallex') return METHOD_ICONS.airwallex
  if (type === 'payoneer') return METHOD_ICONS.stripe
  if (type === 'paypal') return METHOD_ICONS.stripe
  return METHOD_ICONS[type] || alipayIcon
}

function methodSelectedClass(type: string): string {
  if (type.includes('alipay')) return 'border-[#02A9F1] bg-blue-50 text-gray-900 shadow-sm dark:bg-blue-950 dark:text-gray-100'
  if (type.includes('wxpay')) return 'border-[#09BB07] bg-green-50 text-gray-900 shadow-sm dark:bg-green-950 dark:text-gray-100'
  if (type === 'stripe') return 'border-[#676BE5] bg-indigo-50 text-gray-900 shadow-sm dark:bg-indigo-950 dark:text-gray-100'
  if (type === 'airwallex') return 'border-[#FF6B3D] bg-orange-50 text-gray-900 shadow-sm dark:border-[#FF8E3C] dark:bg-orange-950 dark:text-gray-100'
  if (type === 'payoneer') return 'border-[#FF4800] bg-orange-50 text-gray-900 shadow-sm dark:bg-orange-950 dark:text-gray-100'
  if (type === 'paypal') return 'border-[#0070BA] bg-sky-50 text-gray-900 shadow-sm dark:bg-sky-950 dark:text-gray-100'
  return 'border-primary-500 bg-primary-50 text-gray-900 shadow-sm dark:bg-primary-950 dark:text-gray-100'
}
</script>
