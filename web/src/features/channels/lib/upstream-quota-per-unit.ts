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
 * Quota points per currency unit that New API uses by default. Doubao video
 * channels backed by a New API compatible supplier divide the key's remaining
 * quota by this to get the channel balance in the supplier's currency unit.
 */
export const DEFAULT_UPSTREAM_QUOTA_PER_UNIT = 500000

export function isValidUpstreamQuotaPerUnit(value: unknown): value is number {
  return typeof value === 'number' && Number.isFinite(value) && value > 0
}

/**
 * Reads the configured quota per unit from channel settings, given as the
 * stored JSON string or an already parsed object. Falls back to the default.
 */
export function getUpstreamQuotaPerUnit(settings: unknown): number {
  let parsed: unknown = settings
  if (typeof settings === 'string') {
    if (!settings.trim()) return DEFAULT_UPSTREAM_QUOTA_PER_UNIT
    try {
      parsed = JSON.parse(settings)
    } catch {
      return DEFAULT_UPSTREAM_QUOTA_PER_UNIT
    }
  }
  if (parsed && typeof parsed === 'object') {
    const value = (parsed as Record<string, unknown>).upstream_quota_per_unit
    if (isValidUpstreamQuotaPerUnit(value)) return value
  }
  return DEFAULT_UPSTREAM_QUOTA_PER_UNIT
}
