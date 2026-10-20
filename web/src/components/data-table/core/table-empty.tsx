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
import { useTranslation } from 'react-i18next'

import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty'
import { TableRow, TableCell } from '@/components/ui/table'
import { cn } from '@/lib/utils'

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
/** The empty-state block itself (icon, title, description, action). */
export function TableEmptyContent({
  title,
  description,
  icon,
  children,
}: Omit<TableEmptyProps, 'colSpan'>) {
  const { t } = useTranslation()
  const resolvedTitle = title ?? t('No Data')
  const resolvedDescription =
    description ?? t('No records found. Try adjusting your filters.')
  return (
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
  )
}

/**
 * Body spacer for an empty table: reserves the height while the visible
 * empty state is drawn by {@link TableEmptyOverlay} outside the <table>.
 */
export function TableEmptySpacer(props: {
  colSpan: number
  className?: string
}) {
  return (
    <TableRow aria-hidden='true' className='hover:bg-transparent'>
      <TableCell
        colSpan={props.colSpan}
        className={cn('h-[400px] p-0', props.className)}
      />
    </TableRow>
  )
}

/**
 * Empty state drawn as an absolute layer in the table's NON-scrolling outer
 * wrapper (which must be `relative`), below the header. Because it is not
 * inside the horizontally scrolling area it stays centred in the visible
 * box however wide the table is or wherever it is scrolled.
 */
export function TableEmptyOverlay(props: {
  children: React.ReactNode
  className?: string
}) {
  return (
    <div
      data-slot='table-empty-overlay'
      className={cn(
        'pointer-events-none absolute inset-x-0 bottom-0 top-(--table-empty-top,2.5rem) grid place-items-center p-4 [&>*]:pointer-events-auto',
        props.className
      )}
    >
      {props.children}
    </div>
  )
}

/**
 * @deprecated kept for callers that render rows themselves; draws the empty
 * state inside the table (moves with horizontal scroll).
 */
export function TableEmpty({ colSpan, ...rest }: TableEmptyProps) {
  return (
    <TableRow>
      <TableCell colSpan={colSpan} className='h-[400px] p-0'>
        <TableEmptyContent {...rest} />
      </TableCell>
    </TableRow>
  )
}
