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
import i18n, { type Resource } from 'i18next'
import { initReactI18next } from 'react-i18next'

import {
  DEFAULT_INTERFACE_LANGUAGE,
  SUPPORTED_INTERFACE_LANGUAGES,
  type InterfaceLanguageCode,
  detectInterfaceLanguage,
  getBrowserLanguages,
} from './languages'
import en from './locales/en.json'

type LocaleModule = { default: Resource[string] }

/**
 * Only English (the fallback) is bundled with the entry chunk. Every other
 * locale is a separate async chunk fetched once, before the first render, so
 * the page never flashes English text and visitors download one language.
 */
const LOCALE_LOADERS: Record<
  Exclude<InterfaceLanguageCode, 'en'>,
  () => Promise<LocaleModule>
> = {
  zhCN: () => import('./locales/zh.json'),
  zhTW: () => import('./locales/zh-TW.json'),
  fr: () => import('./locales/fr.json'),
  ru: () => import('./locales/ru.json'),
  ja: () => import('./locales/ja.json'),
  vi: () => import('./locales/vi.json'),
}

// The interface language always follows the browser/system language; there
// is no manual switcher and nothing is cached. Clear the value the former
// language detector persisted so it cannot pin an old choice.
const LEGACY_LANGUAGE_STORAGE_KEY = 'i18nextLng'
try {
  if (typeof window !== 'undefined') {
    window.localStorage.removeItem(LEGACY_LANGUAGE_STORAGE_KEY)
  }
} catch {
  /* Storage can be unavailable in private mode. */
}

const detected = detectInterfaceLanguage(getBrowserLanguages())

/** Make sure a locale's resources are registered with i18next. */
export async function loadInterfaceLanguage(
  lng: InterfaceLanguageCode
): Promise<void> {
  if (lng === 'en' || i18n.hasResourceBundle(lng, 'translation')) return
  const mod = await LOCALE_LOADERS[lng]()
  const translation = (mod.default as { translation?: object }).translation
  i18n.addResourceBundle(lng, 'translation', translation ?? {}, true, true)
}

void i18n.use(initReactI18next).init({
  resources: { en },
  lng: detected,
  fallbackLng: DEFAULT_INTERFACE_LANGUAGE,
  supportedLngs: [...SUPPORTED_INTERFACE_LANGUAGES],
  load: 'currentOnly',
  partialBundledLanguages: true,
  nsSeparator: false, // Allow literal colons in keys (e.g., URLs, labels)
  debug: import.meta.env.DEV,
  interpolation: {
    escapeValue: false, // not needed for react as it escapes by default
  },
})

/**
 * Resolves once the detected language is ready. The app waits for this
 * before mounting; a failed chunk load falls back to English instead of
 * blocking the page.
 */
export const i18nReady: Promise<void> = loadInterfaceLanguage(detected)
  .then(() => i18n.changeLanguage(detected))
  .then(
    () => undefined,
    () => undefined
  )

export default i18n
