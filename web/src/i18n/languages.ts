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
/** Interface language codes, as used by i18next resources. */
export const SUPPORTED_INTERFACE_LANGUAGES = [
  'en',
  'zhCN',
  'zhTW',
  'fr',
  'ru',
  'ja',
  'vi',
] as const

export type InterfaceLanguageCode =
  (typeof SUPPORTED_INTERFACE_LANGUAGES)[number]

export const DEFAULT_INTERFACE_LANGUAGE: InterfaceLanguageCode = 'en'

const TRADITIONAL_CHINESE_REGIONS = new Set(['tw', 'hk', 'mo'])

/**
 * Map one BCP-47 tag (as reported by the browser) onto an interface language,
 * or `undefined` when it is not supported.
 *
 * Chinese: `zh-Hant*`, `zh-TW`, `zh-HK`, `zh-MO` -> `zhTW`; every other `zh*`
 * tag (`zh`, `zh-CN`, `zh-SG`, `zh-Hans*`) -> `zhCN`.
 */
export function mapBrowserLanguage(
  tag: string
): InterfaceLanguageCode | undefined {
  const subtags = tag.trim().replaceAll('_', '-').toLowerCase().split('-')
  const primary = subtags[0]
  if (!primary) return undefined
  if (primary === 'zh') {
    if (subtags.includes('hant')) return 'zhTW'
    if (subtags.includes('hans')) return 'zhCN'
    return subtags.some((s) => TRADITIONAL_CHINESE_REGIONS.has(s))
      ? 'zhTW'
      : 'zhCN'
  }
  return SUPPORTED_INTERFACE_LANGUAGES.find(
    (code): code is InterfaceLanguageCode => code === primary
  )
}

/**
 * Pick the interface language from the browser's preferred languages, which
 * follow the browser or operating-system language setting.
 *
 * Only the first (primary) preference counts: if it is not supported the
 * interface falls back to English rather than to a lower-ranked preference.
 */
export function detectInterfaceLanguage(
  languages: readonly string[] | undefined
): InterfaceLanguageCode {
  const primary = languages?.find((tag) => tag && tag.trim() !== '')
  if (!primary) return DEFAULT_INTERFACE_LANGUAGE
  return mapBrowserLanguage(primary) ?? DEFAULT_INTERFACE_LANGUAGE
}

/** The browser's preferred languages, primary first. */
export function getBrowserLanguages(): readonly string[] {
  if (typeof navigator === 'undefined') return []
  if (navigator.languages && navigator.languages.length > 0) {
    return navigator.languages
  }
  return navigator.language ? [navigator.language] : []
}

/**
 * Convert an interface language code (the values i18next uses, such as `zhCN` /
 * `zhTW`) into a valid BCP-47 locale tag that the `Intl.*` APIs accept.
 *
 * `new Intl.NumberFormat('zhCN')` throws `RangeError: Invalid language tag`, so
 * any locale derived from `i18n.language` / `i18n.resolvedLanguage` MUST be run
 * through this before it reaches an `Intl` constructor. Unknown values fall back
 * to `undefined`, which makes `Intl` use the runtime default locale.
 */
export function toIntlLocale(value?: string | null): string | undefined {
  if (!value) return undefined
  switch (value) {
    case 'zhCN':
      return 'zh-CN'
    case 'zhTW':
      return 'zh-TW'
    default:
      break
  }
  try {
    return Intl.getCanonicalLocales(value)[0]
  } catch {
    return undefined
  }
}
