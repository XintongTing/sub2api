<template>
  <header class="sticky top-0 z-40 border-b border-slate-200 bg-white/95 backdrop-blur dark:border-dark-800 dark:bg-dark-950/95">
    <nav class="mx-auto flex h-16 max-w-7xl items-center gap-2 px-3 sm:gap-4 sm:px-6 lg:px-8">
      <router-link to="/home" class="flex min-w-0 shrink-0 items-center gap-2 sm:gap-3">
        <span class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-primary-50 shadow-sm shadow-primary-600/10">
          <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
        </span>
        <span class="hidden max-w-[9rem] truncate text-base font-bold text-slate-950 dark:text-white sm:inline lg:max-w-none">{{ siteName }}</span>
      </router-link>

      <div class="flex min-w-0 flex-1 items-center gap-1 overflow-x-auto px-1 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="shrink-0 whitespace-nowrap rounded-md px-2 py-2 text-xs font-semibold text-slate-700 transition-colors hover:bg-slate-100 hover:text-slate-950 dark:text-dark-200 dark:hover:bg-dark-800 dark:hover:text-white sm:px-3 sm:text-sm"
          :class="{ 'bg-primary-50 text-primary-700 dark:bg-primary-500/10 dark:text-primary-300': isActive(item.path) }"
        >
          {{ item.label }}
        </router-link>
      </div>

      <div class="flex shrink-0 items-center gap-1 sm:gap-2">
        <LocaleSwitcher />
        <button
          type="button"
          class="hidden rounded-lg p-2 text-slate-500 transition-colors hover:bg-slate-100 hover:text-slate-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white sm:inline-flex"
          :title="isDark ? t('publicNav.lightMode') : t('publicNav.darkMode')"
          @click="toggleTheme"
        >
          <Icon v-if="isDark" name="sun" size="md" />
          <Icon v-else name="moon" size="md" />
        </button>
        <router-link
          :to="isAuthenticated ? '/dashboard' : '/login'"
          class="hidden rounded-md bg-primary-600 px-4 py-2 text-sm font-semibold text-white shadow-sm shadow-primary-600/20 transition-colors hover:bg-primary-700 sm:inline-flex"
        >
          {{ isAuthenticated ? t('publicNav.console') : t('publicNav.login') }}
        </router-link>
        <button
          type="button"
          class="rounded-lg p-2 text-slate-600 dark:text-dark-200 sm:hidden"
          :aria-label="t('publicNav.toggleNavigation')"
          @click="mobileOpen = !mobileOpen"
        >
          <Icon name="menu" size="md" />
        </button>
      </div>
    </nav>

    <div v-if="mobileOpen" class="border-t border-slate-200 bg-white px-4 py-3 dark:border-dark-800 dark:bg-dark-950 sm:hidden">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        class="block rounded-md px-3 py-2 text-sm font-semibold text-slate-700 dark:text-dark-200"
        :class="{ 'bg-primary-50 text-primary-700 dark:bg-primary-500/10 dark:text-primary-300': isActive(item.path) }"
        @click="mobileOpen = false"
      >
        {{ item.label }}
      </router-link>
      <router-link
        :to="isAuthenticated ? '/dashboard' : '/login'"
        class="mt-2 block rounded-md bg-primary-600 px-3 py-2 text-center text-sm font-semibold text-white"
        @click="mobileOpen = false"
      >
        {{ isAuthenticated ? t('publicNav.console') : t('publicNav.login') }}
      </router-link>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const route = useRoute()
const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const mobileOpen = ref(false)
const isDark = ref(document.documentElement.classList.contains('dark'))

const siteName = computed(() => appStore.siteName || 'OneAPI')
const siteLogo = computed(() => appStore.siteLogo)
const isAuthenticated = computed(() => authStore.isAuthenticated)

const navItems = computed(() => [
  { path: '/home', label: t('publicNav.home') },
  { path: '/dashboard', label: t('publicNav.console') },
  { path: '/models', label: t('publicNav.models') },
  { path: '/docs', label: t('publicNav.docs') },
])

function isActive(path: string): boolean {
  if (path === '/home') return route.path === '/' || route.path === '/home'
  return route.path === path || route.path.startsWith(`${path}/`)
}

function toggleTheme(): void {
  const nextDark = !document.documentElement.classList.contains('dark')
  document.documentElement.classList.toggle('dark', nextDark)
  localStorage.setItem('theme', nextDark ? 'dark' : 'light')
  isDark.value = nextDark
}
</script>
