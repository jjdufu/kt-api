import { Database } from 'lucide-react'
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
import * as React from 'react'
import { useTranslation } from 'react-i18next'

import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty'
import { TableRow, TableCell } from '@/components/ui/table'

interface TableEmptyProps {
  /**
   * Number of columns to span
   */
  colSpan: number
  /**
   * Custom title for empty state
   * @default 'No Data'
   */
  title?: string
  /**
   * Custom description for empty state
   * @default 'No records found. Try adjusting your filters.'
   */
  description?: string
  /**
   * Custom icon component
   * @default Database icon
   */
  icon?: React.ReactNode
  /**
   * Additional content to display (e.g., buttons)
   */
  children?: React.ReactNode
}

/**
 * Generic table empty state component
 * Displays a centered empty state message when table has no data
 */
/**
 * Keeps the empty state centred in the *visible* part of a horizontally
 * scrolling table: sticks to the scroll container's left edge and takes the
 * container's client width (tracked with a ResizeObserver), instead of
 * centring across the full width of a wide table.
 */
function ViewportCentered(props: { children: React.ReactNode }) {
  const ref = React.useRef<HTMLDivElement>(null)
  const [width, setWidth] = React.useState<number | undefined>(undefined)
  React.useLayoutEffect(() => {
    const scroller = ref.current?.closest<HTMLElement>(
      "[data-slot='table-container']"
    )
    if (!scroller) return
    const update = () => setWidth(scroller.clientWidth)
    update()
    if (typeof ResizeObserver === 'undefined') return
    const observer = new ResizeObserver(update)
    observer.observe(scroller)
    return () => observer.disconnect()
  }, [])
  return (
    <div
      ref={ref}
      data-slot='table-empty-viewport'
      className='sticky left-0 flex h-full items-center justify-center'
      style={width ? { width } : undefined}
    >
      {props.children}
    </div>
  )
}

export function TableEmpty({
  colSpan,
  title,
  description,
  icon,
  children,
}: TableEmptyProps) {
  const { t } = useTranslation()
  const resolvedTitle = title ?? t('No Data')
  const resolvedDescription =
    description ?? t('No records found. Try adjusting your filters.')
  return (
    <TableRow>
      <TableCell colSpan={colSpan} className='h-[400px] p-0'>
        <ViewportCentered>
          <Empty>
            <EmptyHeader>
              <EmptyMedia variant='icon'>
                {icon || <Database className='size-6' />}
              </EmptyMedia>
              <EmptyTitle>{resolvedTitle}</EmptyTitle>
              <EmptyDescription>{resolvedDescription}</EmptyDescription>
            </EmptyHeader>
            {children}
          </Empty>
        </ViewportCentered>
      </TableCell>
    </TableRow>
  )
}
