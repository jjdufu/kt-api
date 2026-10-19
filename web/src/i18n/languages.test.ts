/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { describe, expect, test } from 'vitest'

import { detectInterfaceLanguage, mapBrowserLanguage } from './languages'

describe('mapBrowserLanguage', () => {
  test('maps Chinese tags by script first, then region', () => {
    for (const tag of [
      'zh-TW',
      'zh-HK',
      'zh-MO',
      'zh-Hant',
      'zh-Hant-TW',
      'zh-hant-cn',
      'zh_TW',
    ]) {
      expect(mapBrowserLanguage(tag)).toBe('zhTW')
    }
    for (const tag of [
      'zh',
      'zh-CN',
      'zh-SG',
      'zh-Hans',
      'zh-Hans-HK',
      'zh-hans-tw',
      'zh_CN',
    ]) {
      expect(mapBrowserLanguage(tag)).toBe('zhCN')
    }
  })

  test('maps other supported languages by primary subtag', () => {
    expect(mapBrowserLanguage('en-US')).toBe('en')
    expect(mapBrowserLanguage('fr-FR')).toBe('fr')
    expect(mapBrowserLanguage('ru')).toBe('ru')
    expect(mapBrowserLanguage('ja-JP')).toBe('ja')
    expect(mapBrowserLanguage('vi-VN')).toBe('vi')
  })

  test('rejects unsupported languages', () => {
    expect(mapBrowserLanguage('de-DE')).toBe(undefined)
    expect(mapBrowserLanguage('zhCN')).toBe(undefined)
    expect(mapBrowserLanguage('')).toBe(undefined)
  })
})

describe('detectInterfaceLanguage', () => {
  test('uses only the primary browser language', () => {
    expect(detectInterfaceLanguage(['zh-CN', 'en'])).toBe('zhCN')
    expect(detectInterfaceLanguage(['ja', 'zh-CN'])).toBe('ja')
    // An unsupported primary language falls back to English, not to zh-CN.
    expect(detectInterfaceLanguage(['de-DE', 'zh-CN'])).toBe('en')
  })

  test('falls back to English without a preference', () => {
    expect(detectInterfaceLanguage([])).toBe('en')
    expect(detectInterfaceLanguage(undefined)).toBe('en')
  })
})
