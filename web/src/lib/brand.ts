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
import i18n from 'i18next'

import { DEFAULT_SYSTEM_NAME } from './constants'

/**
 * The built-in site name is bilingual: 开天大模型 in Chinese locales and
 * KTAPI elsewhere (see the "开天大模型" locale key). A name configured by an
 * admin is shown exactly as entered.
 */
export function localizeSystemName(name: string | undefined | null): string {
  const value = (name ?? '').trim() || DEFAULT_SYSTEM_NAME
  if (value !== DEFAULT_SYSTEM_NAME) return value
  return i18n.t(DEFAULT_SYSTEM_NAME, { defaultValue: DEFAULT_SYSTEM_NAME })
}
