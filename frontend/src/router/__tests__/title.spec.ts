import { describe, expect, it } from 'vitest'
import { resolveDocumentTitle } from '@/router/title'

describe('resolveDocumentTitle', () => {
  it('uses route title and site name when both are present', () => {
    expect(resolveDocumentTitle('Usage Records', 'My Site')).toBe('Usage Records - My Site')
  })

  it('falls back to site name when route title is missing', () => {
    expect(resolveDocumentTitle(undefined, 'My Site')).toBe('My Site')
  })

  it('falls back to OneAPI when site name is blank', () => {
    expect(resolveDocumentTitle('Dashboard', '')).toBe('Dashboard - OneAPI')
    expect(resolveDocumentTitle(undefined, '   ')).toBe('OneAPI')
  })

  it('uses the latest site name without persisting stale titles', () => {
    const before = resolveDocumentTitle('Admin Dashboard', 'Alpha')
    const after = resolveDocumentTitle('Admin Dashboard', 'Beta')

    expect(before).toBe('Admin Dashboard - Alpha')
    expect(after).toBe('Admin Dashboard - Beta')
  })
})
