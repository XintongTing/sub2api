import { describe, expect, it } from 'vitest'
import { availableLocales, normalizeLocaleCode } from '@/i18n'

describe('OneAPI locale routing', () => {
  it('exposes English, Thai, Simplified Chinese, and Traditional Chinese', () => {
    expect(availableLocales.map(locale => locale.code)).toEqual(['en', 'th', 'zh-CN', 'zh-TW'])
  })

  it('normalizes legacy and browser locale codes', () => {
    expect(normalizeLocaleCode('zh')).toBe('zh-CN')
    expect(normalizeLocaleCode('zh_Hans')).toBe('zh-CN')
    expect(normalizeLocaleCode('zh-HK')).toBe('zh-TW')
    expect(normalizeLocaleCode('th-TH')).toBe('th')
  })

  it('ships localized public home copy for Thai and Traditional Chinese', async () => {
    const [{ default: th }, { default: zhTW }] = await Promise.all([
      import('../locales/th'),
      import('../locales/zh-TW'),
    ])

    expect(th.home.defaultSiteSubtitle).toBe('เกตเวย์ API สำหรับเรียกใช้โมเดล AI หลัก')
    expect(th.home.features.unifiedGateway).toBe('เชื่อมต่อในคลิกเดียว')
    expect(th.home.terminal.routingComment).toBe('# กำลังส่งต่อไปยัง upstream...')
    expect(zhTW.home.defaultSiteSubtitle).toBe('AI 模型 API 聚合網關')
    expect(zhTW.home.features.unifiedGateway).toBe('一鍵接入')
    expect(zhTW.home.tags.stickySession).toBe('會話保持')
  })
})
