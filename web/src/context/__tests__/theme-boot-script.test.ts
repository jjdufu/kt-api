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
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { THEME_STORAGE_KEYS } from '@/lib/theme-storage'

// The inline script in index.html must apply the theme class before React
// mounts; Dark Reader checks `html.dark` only once while the page loads.
const html = readFileSync(
  resolve(import.meta.dirname, '../../../index.html'),
  'utf8'
)
const scripts = [...html.matchAll(/<script>([\s\S]*?)<\/script>/g)].map(
  (m) => m[1]
)
const bootScript = scripts[0] ?? ''
const presetScript = scripts[1] ?? ''

function runBootScript() {
  new Function(bootScript)()
}

function mockSystemDarkMode(dark: boolean) {
  vi.spyOn(window, 'matchMedia').mockImplementation(
    (query: string) =>
      ({
        matches: dark && query === '(prefers-color-scheme: dark)',
        media: query,
      }) as MediaQueryList
  )
}

beforeEach(() => {
  localStorage.clear()
})

afterEach(() => {
  vi.restoreAllMocks()
  localStorage.clear()
  document.documentElement.classList.remove('light', 'dark')
})

describe('index.html theme boot script', () => {
  it('applies the saved dark preference to html before React mounts', () => {
    mockSystemDarkMode(false)
    localStorage.setItem(THEME_STORAGE_KEYS.mode, 'dark')

    runBootScript()

    expect(document.documentElement).toHaveClass('dark')
    expect(document.documentElement).not.toHaveClass('light')
  })

  it('applies the saved light preference even when the system is dark', () => {
    mockSystemDarkMode(true)
    localStorage.setItem(THEME_STORAGE_KEYS.mode, 'light')

    runBootScript()

    expect(document.documentElement).toHaveClass('light')
    expect(document.documentElement).not.toHaveClass('dark')
  })

  it('follows the system dark mode when no preference is saved', () => {
    mockSystemDarkMode(true)

    runBootScript()

    expect(document.documentElement).toHaveClass('dark')
  })

  it('follows the system light mode when the saved value is invalid', () => {
    mockSystemDarkMode(false)
    localStorage.setItem(THEME_STORAGE_KEYS.mode, 'sepia')

    runBootScript()

    expect(document.documentElement).toHaveClass('light')
  })
})

describe('index.html preset boot script', () => {
  afterEach(() => document.body.removeAttribute('data-theme-preset'))

  it('applies KT Comic for fresh visitors without storing it', () => {
    new Function(presetScript)()
    expect(document.body).toHaveAttribute('data-theme-preset', 'kt-comic')
    expect(localStorage.getItem(THEME_STORAGE_KEYS.preset)).toBeNull()
  })

  afterEach(() => document.body.removeAttribute('data-theme-large-font'))

  it.each(['default', 'kt-comic'])(
    'keeps supported stored preset %s',
    (preset) => {
      localStorage.setItem(THEME_STORAGE_KEYS.preset, preset)
      new Function(presetScript)()
      expect(document.body).toHaveAttribute('data-theme-preset', preset)
      expect(document.body).not.toHaveAttribute('data-theme-large-font')
    }
  )

  it('applies a stored large-font preference with any preset', () => {
    localStorage.setItem(THEME_STORAGE_KEYS.largeFont, '1')
    new Function(presetScript)()
    expect(document.body).toHaveAttribute('data-theme-preset', 'kt-comic')
    expect(document.body).toHaveAttribute('data-theme-large-font', '')
  })

  it.each(['simple', 'bogus'])('ignores unknown preset %s', (value) => {
    localStorage.setItem(THEME_STORAGE_KEYS.preset, value)
    new Function(presetScript)()
    expect(document.body).toHaveAttribute('data-theme-preset', 'kt-comic')
  })
})
