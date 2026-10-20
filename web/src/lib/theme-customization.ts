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
/**
 * Theme customization constants and types.
 *
 * Lives in `lib/` (not `context/`) so it can be imported alongside the
 * provider without breaking React Fast Refresh boundaries.
 */

export const THEME_PRESETS = [
  {
    // 开天漫画: the site's comic / manga cyberpunk look (paper & ink light,
    // neon-on-ink dark). Site default for visitors without a saved choice.
    value: 'kt-comic',
    name: 'KT Comic',
    swatches: ['#ffe600', '#ff2bd6'],
  },
  {
    // 暗夜: muted green/magenta night palette (light & dark).
    value: 'underground',
    name: 'Underground',
    swatches: ['oklch(0.5315 0.0694 156.19)', 'oklch(0.5748 0.0862 336.52)'],
  },
  {
    // 超大字体简易: high-contrast neutral surfaces with larger type.
    value: 'simple-large',
    name: 'Simple Large-font',
    swatches: ['oklch(0.15 0 0)', 'oklch(0.99 0 0)'],
  },
] as const

export type SelectableThemePreset = (typeof THEME_PRESETS)[number]['value']
/**
 * `default` is the hidden neutral fallback (theme.css palette, applied by
 * removing `data-theme-preset`). It is never offered in the picker and a
 * stored `default` is migrated like any other removed preset.
 */
export type ThemePreset = SelectableThemePreset | 'default'
export type ThemeRadius = 'default' | 'none' | 'sm' | 'md' | 'lg' | 'xl'
export type ThemeScale = 'default' | 'sm' | 'lg' | 'xl'
export type ContentLayout = 'full' | 'centered'

/**
 * Font axis for the theme.
 *
 * - `default` — use the browser's sans-serif unless the active preset
 *   specifies its own font.
 * - `sans` — explicitly use Public Sans.
 * - `serif` — use Lora + CJK fallbacks. Monospace contexts keep their
 *   own family via Tailwind preflight and `.font-mono`.
 */
export type ThemeFont = 'default' | 'sans' | 'serif'

/**
 * The resolved (non-`default`) font value applied to the DOM. The provider
 * always sets `data-theme-font` to one of these concrete values so CSS only
 * needs simple attribute selectors (no `:not()` gymnastics, no per-preset
 * font branches).
 */
export type ResolvedThemeFont = 'system' | Exclude<ThemeFont, 'default'>

export type ThemeCustomization = {
  preset: ThemePreset
  font: ThemeFont
  radius: ThemeRadius
  scale: ThemeScale
  contentLayout: ContentLayout
}

/**
 * The neutral palette from theme.css. It is applied by *removing*
 * `data-theme-preset`. Internal only; not selectable.
 */
export const NEUTRAL_THEME_PRESET: ThemePreset = 'default'

export const DEFAULT_THEME_CUSTOMIZATION: ThemeCustomization = {
  // New visitors (no stored preference) get the KT Comic look. Users can
  // still pick any other preset; that choice is persisted in localStorage.
  preset: 'kt-comic',
  font: 'default',
  radius: 'default',
  scale: 'default',
  contentLayout: 'full',
}

export const THEME_PRESET_VALUES = new Set(
  THEME_PRESETS.map((p) => p.value)
) as ReadonlySet<ThemePreset>

/**
 * Preset a stored-but-removed value (any former preset, including the old
 * `default`) is migrated to. Fresh visitors without a stored value still get
 * DEFAULT_THEME_CUSTOMIZATION.preset. Keep in sync with index.html.
 */
export const MIGRATED_THEME_PRESET: ThemePreset = 'underground'

export function migrateThemePreset(raw: string | null): ThemePreset | null {
  if (raw === null || raw === '') return null
  if (THEME_PRESET_VALUES.has(raw as ThemePreset)) return raw as ThemePreset
  return MIGRATED_THEME_PRESET
}

export const THEME_FONT_VALUES: ReadonlySet<ThemeFont> = new Set([
  'default',
  'sans',
  'serif',
])

export const THEME_RADIUS_VALUES: ReadonlySet<ThemeRadius> = new Set([
  'default',
  'none',
  'sm',
  'md',
  'lg',
  'xl',
])

export const THEME_SCALE_VALUES: ReadonlySet<ThemeScale> = new Set([
  'default',
  'sm',
  'lg',
  'xl',
])

export const CONTENT_LAYOUT_VALUES: ReadonlySet<ContentLayout> = new Set([
  'full',
  'centered',
])

/**
 * Preset → default font mapping. Used by the provider to resolve the user's
 * `font: 'default'` preference against the active preset.
 *
 * Presets without a font override use the browser's sans-serif. Public
 * Sans is reserved for an explicit user choice.
 */
export const PRESET_DEFAULT_FONT: Partial<
  Record<ThemePreset, ResolvedThemeFont>
> = {}

/**
 * Resolve a user font preference + active preset into the concrete font that
 * should drive the DOM. Pure function so it's safe to call inside both the
 * effect that applies the attribute and the UI preview that hints at what
 * `default` will render as.
 */
export function resolveThemeFont(
  font: ThemeFont,
  preset: ThemePreset
): ResolvedThemeFont {
  if (font === 'default') {
    return PRESET_DEFAULT_FONT[preset] ?? 'system'
  }
  return font
}
