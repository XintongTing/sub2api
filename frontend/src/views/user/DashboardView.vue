<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-6">
      <section class="overflow-hidden rounded-lg border border-primary-100 bg-white shadow-sm dark:border-dark-800 dark:bg-dark-900">
        <div class="bg-gradient-to-r from-primary-700 via-primary-600 to-cyan-500 px-6 py-6 text-white">
          <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
            <div>
              <p class="text-sm font-semibold text-primary-50">{{ t('dashboard.consoleOverview') }}</p>
              <h1 class="mt-2 text-2xl font-bold sm:text-3xl">{{ t('dashboard.title') }}</h1>
              <p class="mt-2 text-sm leading-6 text-primary-50">
                {{ t('dashboard.overviewDescription') }}
              </p>
            </div>
            <button
              type="button"
              class="inline-flex items-center justify-center rounded-md bg-white px-4 py-2 text-sm font-bold text-primary-700 transition-colors hover:bg-primary-50 disabled:cursor-not-allowed disabled:opacity-70"
              :disabled="isRefreshing"
              @click="refreshAll"
            >
              {{ isRefreshing ? t('dashboard.refreshing') : t('dashboard.refreshData') }}
            </button>
          </div>
        </div>
      </section>

      <div
        v-if="statsError"
        class="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm leading-6 text-amber-800 dark:border-amber-900/50 dark:bg-amber-950/30 dark:text-amber-200"
      >
        {{ statsError }}
      </div>

      <UserDashboardStats
        :stats="displayStats"
        :balance="currentBalance"
        :is-simple="authStore.isSimpleMode"
        :platform-quotas="platformQuotas"
      />

      <UserDashboardCharts
        v-model:startDate="startDate"
        v-model:endDate="endDate"
        v-model:granularity="granularity"
        :loading="loadingCharts"
        :trend="trendData"
        :models="modelStats"
        @dateRangeChange="loadCharts"
        @granularityChange="loadCharts"
        @refresh="refreshAll"
      />

      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div class="lg:col-span-2">
          <UserDashboardRecentUsage :data="recentUsage" :loading="loadingUsage" />
        </div>
        <div class="lg:col-span-1">
          <UserDashboardQuickActions />
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { usageAPI, type UserDashboardStats as UserStatsType } from '@/api/usage'
import AppLayout from '@/components/layout/AppLayout.vue'
import UserDashboardStats from '@/components/user/dashboard/UserDashboardStats.vue'
import UserDashboardCharts from '@/components/user/dashboard/UserDashboardCharts.vue'
import UserDashboardRecentUsage from '@/components/user/dashboard/UserDashboardRecentUsage.vue'
import UserDashboardQuickActions from '@/components/user/dashboard/UserDashboardQuickActions.vue'
import type { UsageLog, TrendDataPoint, ModelStat, PlatformQuotaItem } from '@/types'

const EMPTY_STATS: UserStatsType = {
  balance: 0,
  total_api_keys: 0,
  active_api_keys: 0,
  total_requests: 0,
  total_input_tokens: 0,
  total_output_tokens: 0,
  total_cache_creation_tokens: 0,
  total_cache_read_tokens: 0,
  total_tokens: 0,
  total_cost: 0,
  total_actual_cost: 0,
  today_requests: 0,
  today_input_tokens: 0,
  today_output_tokens: 0,
  today_cache_creation_tokens: 0,
  today_cache_read_tokens: 0,
  today_tokens: 0,
  today_cost: 0,
  today_actual_cost: 0,
  average_duration_ms: 0,
  rpm: 0,
  tpm: 0,
  by_platform: []
}

const authStore = useAuthStore()
const { t } = useI18n()
const user = computed(() => authStore.user)
const stats = ref<UserStatsType | null>(null)
const statsError = ref('')
const loadingStats = ref(false)
const loadingUsage = ref(false)
const loadingCharts = ref(false)
const loadingPlatformQuotas = ref(false)
const trendData = ref<TrendDataPoint[]>([])
const modelStats = ref<ModelStat[]>([])
const recentUsage = ref<UsageLog[]>([])
const platformQuotas = ref<PlatformQuotaItem[] | null>(null)

const displayStats = computed(() => stats.value ?? EMPTY_STATS)
const currentBalance = computed(() => {
  const statsBalance = stats.value?.balance
  if (typeof statsBalance === 'number' && Number.isFinite(statsBalance)) {
    return statsBalance
  }
  return user.value?.balance ?? 0
})
const isRefreshing = computed(() => loadingStats.value || loadingCharts.value || loadingUsage.value || loadingPlatformQuotas.value)

const formatLD = (d: Date) => d.toISOString().split('T')[0]
const startDate = ref(formatLD(new Date(Date.now() - 6 * 86400000)))
const endDate = ref(formatLD(new Date()))
const granularity = ref<'day' | 'hour'>('day')

function messageFromError(error: unknown): string {
  if (error instanceof Error) return error.message
  if (typeof error === 'object' && error !== null && 'message' in error) {
    return String((error as { message?: unknown }).message || t('dashboard.requestFailed'))
  }
  return t('dashboard.requestFailed')
}

function withTimeout<T>(promise: Promise<T>, ms: number, label: string): Promise<T> {
  return Promise.race([
    promise,
    new Promise<T>((_, reject) => {
      window.setTimeout(() => reject(new Error(t('dashboard.requestTimeout', { label }))), ms)
    })
  ])
}

async function loadStats(): Promise<void> {
  loadingStats.value = true
  statsError.value = ''

  try {
    const [userResult, statsResult] = await Promise.allSettled([
      authStore.refreshUser(),
      withTimeout(usageAPI.getDashboardStats(), 12000, t('dashboard.statsLabel'))
    ])

    if (userResult.status === 'rejected') {
      console.warn('Failed to refresh current user:', userResult.reason)
    }

    if (statsResult.status === 'fulfilled') {
      stats.value = statsResult.value
      if (user.value && typeof statsResult.value.balance === 'number' && Number.isFinite(statsResult.value.balance)) {
        authStore.user = {
          ...user.value,
          balance: statsResult.value.balance
        }
      }
      return
    }

    stats.value = stats.value ?? EMPTY_STATS
    statsError.value = t('dashboard.statsUnavailable', { message: messageFromError(statsResult.reason) })
    console.error('Failed to load dashboard stats:', statsResult.reason)
  } finally {
    loadingStats.value = false
  }
}

async function loadCharts(): Promise<void> {
  loadingCharts.value = true
  try {
    const [trend, models] = await Promise.all([
      usageAPI.getDashboardTrend({
        start_date: startDate.value,
        end_date: endDate.value,
        granularity: granularity.value
      }),
      usageAPI.getDashboardModels({
        start_date: startDate.value,
        end_date: endDate.value
      })
    ])
    trendData.value = trend.trend || []
    modelStats.value = models.models || []
  } catch (error) {
    console.error('Failed to load charts:', error)
    trendData.value = []
    modelStats.value = []
  } finally {
    loadingCharts.value = false
  }
}

async function loadRecent(): Promise<void> {
  loadingUsage.value = true
  try {
    const res = await usageAPI.getByDateRange(startDate.value, endDate.value)
    recentUsage.value = res.items.slice(0, 5)
  } catch (error) {
    console.error('Failed to load recent usage:', error)
    recentUsage.value = []
  } finally {
    loadingUsage.value = false
  }
}

async function loadPlatformQuotas(): Promise<void> {
  platformQuotas.value = []
  loadingPlatformQuotas.value = false
}

function refreshAll(): void {
  void loadStats()
  void loadCharts()
  void loadRecent()
  void loadPlatformQuotas()
}

onMounted(() => {
  refreshAll()
})
</script>
