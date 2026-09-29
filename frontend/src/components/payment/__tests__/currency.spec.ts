import { describe, expect, it } from 'vitest'
import { formatPaymentAmount, normalizePaymentCurrency } from '../currency'

describe('formatPaymentAmount', () => {
  it('uses the currency default fraction digits', () => {
    expect(formatPaymentAmount(100, 'JPY', 'en-US')).not.toContain('.00')
    expect(formatPaymentAmount(100, 'KRW', 'en-US')).not.toContain('.00')
    expect(formatPaymentAmount(100, 'HKD', 'en-US')).toContain('.00')
  })

  it('defaults invalid or missing currencies to THB', () => {
    expect(normalizePaymentCurrency('')).toBe('THB')
    expect(normalizePaymentCurrency(null)).toBe('THB')
  })
})
