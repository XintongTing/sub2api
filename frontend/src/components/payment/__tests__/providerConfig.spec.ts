import { describe, expect, it } from 'vitest'
import { PAYMENT_CURRENCY_OPTIONS, PROVIDER_CALLBACK_PATHS, PROVIDER_CONFIG_FIELDS, PROVIDER_SUPPORTED_TYPES, WEBHOOK_PATHS } from '@/components/payment/providerConfig'

function findField(providerKey: string, key: string) {
  const fields = PROVIDER_CONFIG_FIELDS[providerKey] || []
  return fields.find(field => field.key === key)
}

describe('PROVIDER_CONFIG_FIELDS.wxpay', () => {
  it('keeps admin form validation aligned with backend-required credentials', () => {
    expect(findField('wxpay', 'publicKeyId')?.optional).toBeFalsy()
    expect(findField('wxpay', 'certSerial')?.optional).toBeFalsy()
  })

  it('only keeps the simplified visible credential set in the admin form', () => {
    expect(findField('wxpay', 'mpAppId')).toBeUndefined()
    expect(findField('wxpay', 'h5AppName')).toBeUndefined()
    expect(findField('wxpay', 'h5AppUrl')).toBeUndefined()
  })
})

describe('PROVIDER_CONFIG_FIELDS.airwallex', () => {
  it('adds currency config with THB as the default', () => {
    const currency = findField('airwallex', 'currency')

    expect(currency?.defaultValue).toBe('THB')
    expect(currency?.hintKey).toBe('admin.settings.payment.field_paymentCurrencyHint')
    expect(currency?.options).toBe(PAYMENT_CURRENCY_OPTIONS)
  })

  it('marks accountId as optional and explains when it can be left blank', () => {
    const accountId = findField('airwallex', 'accountId')

    expect(accountId?.optional).toBe(true)
    expect(accountId?.clearable).toBe(true)
    expect(accountId?.hintKey).toBe('admin.settings.payment.field_accountIdHint')
  })

  it('explains that apiBase must match the Airwallex key environment', () => {
    expect(findField('airwallex', 'apiBase')?.hintKey).toBe('admin.settings.payment.field_airwallexApiBaseHint')
  })
})

describe('PROVIDER_CONFIG_FIELDS.stripe', () => {
  it('adds currency config with THB as the default', () => {
    const currency = findField('stripe', 'currency')

    expect(currency?.defaultValue).toBe('THB')
    expect(currency?.hintKey).toBe('admin.settings.payment.field_paymentCurrencyHint')
    expect(currency?.options).toBe(PAYMENT_CURRENCY_OPTIONS)
  })
})

describe('PROVIDER_CONFIG_FIELDS.payoneer', () => {
  it('defines Payoneer Checkout credentials and THB currency defaults', () => {
    expect(PROVIDER_SUPPORTED_TYPES.payoneer).toEqual(['payoneer'])
    expect(WEBHOOK_PATHS.payoneer).toBe('/api/v1/payment/webhook/payoneer')
    expect(PROVIDER_CALLBACK_PATHS.payoneer?.returnUrl).toBe('/payment/result')
    expect(findField('payoneer', 'clientId')?.sensitive).toBe(false)
    expect(findField('payoneer', 'clientSecret')?.sensitive).toBe(true)
    expect(findField('payoneer', 'webhookSecret')?.sensitive).toBe(true)
    expect(findField('payoneer', 'environment')?.defaultValue).toBe('sandbox')
    expect(findField('payoneer', 'currency')?.defaultValue).toBe('THB')
  })
})
