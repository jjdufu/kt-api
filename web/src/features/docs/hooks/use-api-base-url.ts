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
import { useMemo } from 'react'

import { useStatus } from '@/hooks/use-status'

/**
 * Resolves the API base address shown in the docs: the configured server
 * address (System Settings -> server address) without a trailing slash, else
 * the given origin.
 */
export function resolveApiBaseUrl(serverAddress: unknown, origin: string) {
  if (typeof serverAddress === 'string' && serverAddress.trim()) {
    return serverAddress.trim().replace(/\/+$/, '')
  }
  return origin
}

export function useApiBaseUrl(): string {
  const { status } = useStatus()
  return useMemo(() => {
    const record = status as Record<string, unknown> | null
    const origin =
      typeof window !== 'undefined' && window.location
        ? window.location.origin
        : ''
    return resolveApiBaseUrl(record?.server_address, origin)
  }, [status])
}
