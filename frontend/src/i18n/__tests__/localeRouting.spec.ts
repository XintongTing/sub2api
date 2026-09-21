import { describe, expect, it } from 'vitest'
import { availableLocales, normalizeLocaleCode } from '@/i18n'

describe('OneAPI locale routing', () => {
  it('exposes English, Thai, Simplified Chinese, and Traditional Chinese', () => {
    expect(availableLocales.map(locale => locale.code)).toEqual(['th', 'en', 'zh-TW', 'zh-CN'])
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

    expect(th.home.defaultSiteSubtitle).toBe('แพลตฟอร์ม API สำหรับเชื่อมต่อโมเดล AI หลักด้วย API Key เดียว')
    expect(th.home.features.unifiedGateway).toBe('การเข้าถึงด้วยคลิกเดียว')
    expect(th.home.terminal.routingComment).toBe('# กำลังส่งต่อไปยังต้นทาง...')
    expect(zhTW.home.defaultSiteSubtitle).toBe('AI 模型 API 聚合網關')
    expect(zhTW.home.features.unifiedGateway).toBe('一鍵接入')
    expect(zhTW.home.tags.stickySession).toBe('會話保持')
  })

  it('keeps Thai aligned with the complete English message tree', async () => {
    const [{ default: en }, { default: th }] = await Promise.all([
      import('../locales/en'),
      import('../locales/th'),
    ])

    const leafKeys = (messages: Record<string, any>, prefix = ''): string[] =>
      Object.entries(messages).flatMap(([key, value]) => {
        const path = prefix ? `${prefix}.${key}` : key
        return value && typeof value === 'object'
          ? leafKeys(value, path)
          : [path]
      })

    const thaiKeys = leafKeys(th).filter(key => !['nav.recharge', 'user.profile.title', 'user.profile.subtitle'].includes(key))
    expect(thaiKeys.sort()).toEqual(leafKeys(en).sort())
    expect(th.payment.rechargeHelp).toContain('เลือกวิธีชำระเงิน')
  })
})
