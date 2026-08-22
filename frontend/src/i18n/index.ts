import { createI18n } from 'vue-i18n'

export type LocaleCode = 'en' | 'zh-CN' | 'zh-TW' | 'th'

type LocaleMessages = Record<string, any>

const LOCALE_KEY = 'sub2api_locale'
const DEFAULT_LOCALE: LocaleCode = 'th'
const FALLBACK_LOCALE: LocaleCode = 'zh-CN'

const localeLoaders: Record<LocaleCode, () => Promise<{ default: LocaleMessages }>> = {
  en: () => import('./locales/en'),
  'zh-CN': () => import('./locales/zh-CN'),
  'zh-TW': () => import('./locales/zh-TW'),
  th: () => import('./locales/th'),
}

function isLocaleCode(value: string): value is LocaleCode {
  return value === 'en' || value === 'zh-CN' || value === 'zh-TW' || value === 'th'
}

export function normalizeLocaleCode(value?: string | null): LocaleCode | '' {
  const normalized = String(value || '').trim().replace('_', '-')
  const lower = normalized.toLowerCase()

  if (lower === 'zh') return 'zh-CN'
  if (lower === 'zh-cn' || lower === 'zh-hans') return 'zh-CN'
  if (lower === 'zh-tw' || lower === 'zh-hk' || lower === 'zh-mo' || lower === 'zh-hant') return 'zh-TW'
  if (lower === 'th' || lower.startsWith('th-')) return 'th'
  if (lower === 'en' || lower.startsWith('en-')) return 'en'

  return ''
}

function getDefaultLocale(): LocaleCode {
  const saved = localStorage.getItem(LOCALE_KEY)
  const savedLocale = normalizeLocaleCode(saved)
  if (savedLocale) {
    return savedLocale
  }

  for (const lang of navigator.languages || [navigator.language]) {
    const locale = normalizeLocaleCode(lang)
    if (locale) {
      return locale
    }
  }

  return DEFAULT_LOCALE
}

export const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: FALLBACK_LOCALE,
  messages: {},
  warnHtmlMessage: false,
})

const loadedLocales = new Set<LocaleCode>()

export async function loadLocaleMessages(locale: LocaleCode): Promise<void> {
  if (loadedLocales.has(locale)) {
    return
  }

  const loader = localeLoaders[locale]
  const module = await loader()
  i18n.global.setLocaleMessage(locale, module.default)
  loadedLocales.add(locale)
}

export async function initI18n(): Promise<void> {
  const current = getLocale()
  await loadLocaleMessages(current)
  document.documentElement.setAttribute('lang', current)
}

export async function setLocale(locale: string): Promise<void> {
  const normalized = normalizeLocaleCode(locale)
  if (!normalized || !isLocaleCode(normalized)) {
    return
  }

  await loadLocaleMessages(normalized)
  i18n.global.locale.value = normalized
  localStorage.setItem(LOCALE_KEY, normalized)
  document.documentElement.setAttribute('lang', normalized)

  const { resolveDocumentTitle } = await import('@/router/title')
  const { default: router } = await import('@/router')
  const { useAppStore } = await import('@/stores/app')
  const route = router.currentRoute.value
  const appStore = useAppStore()
  document.title = resolveDocumentTitle(route.meta.title, appStore.siteName, route.meta.titleKey as string)
}

export function getLocale(): LocaleCode {
  const current = i18n.global.locale.value
  return isLocaleCode(current) ? current : DEFAULT_LOCALE
}

export const availableLocales = [
  { code: 'th', name: 'ไทย', flag: 'TH' },
  { code: 'en', name: 'English', flag: 'EN' },
  { code: 'zh-TW', name: '繁體中文', flag: 'TW' },
  { code: 'zh-CN', name: '简体中文', flag: 'CN' },
] as const

export default i18n
