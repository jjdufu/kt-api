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
import type { ComponentProps } from 'react'

import { StatusBadge } from '@/components/status-badge'
import { cn } from '@/lib/utils'

/** Shared unsaved-change treatment; colors, radius and type scale follow the preset. */
export function ChangeBadge({
  className,
  ...props
}: ComponentProps<typeof StatusBadge>) {
  return (
    <StatusBadge
      {...props}
      variant='warning'
      type='badge'
      copyable={false}
      className={cn(
        'shrink-0 rounded-md bg-warning/10 px-2 text-xs text-foreground ring-1 ring-warning/25 ring-inset',
        className
      )}
    />
  )
}
