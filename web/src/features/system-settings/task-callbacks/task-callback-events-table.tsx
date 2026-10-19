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

import { StaticDataTable } from '@/components/data-table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { toIntlLocale } from '@/i18n/languages'
import { formatTimestampRelative, formatTimestampToDate } from '@/lib/format'

import type { TaskCallbackEvent } from './api'
import { taskCallbackStateLabel } from './constants'

const STATE_VARIANT: Record<string, 'secondary' | 'destructive' | 'outline'> = {
  succeeded: 'secondary',
  failed: 'destructive',
  pending: 'outline',
  delivering: 'outline',
}

function lastResult(event: TaskCallbackEvent): string {
  if (event.last_http_status > 0) {
    return `HTTP ${event.last_http_status}`
  }
  return event.last_error || '-'
}

type TaskCallbackEventsTableProps = {
  events: TaskCallbackEvent[]
  showTask?: boolean
  redeliveringId?: string
  onRedeliver?: (event: TaskCallbackEvent) => void
}

/** Delivery events of task callbacks, shared by the admin list and task details. */
export function TaskCallbackEventsTable(props: TaskCallbackEventsTableProps) {
  const { t, i18n } = useTranslation()
  const locale = toIntlLocale(i18n.resolvedLanguage || i18n.language)

  return (
    <StaticDataTable tableClassName='min-w-[860px]'>
      <TableHeader>
        <TableRow className='bg-muted/40 hover:bg-muted/40'>
          {props.showTask && (
            <TableHead className='h-9 px-4 text-xs'>{t('Task ID')}</TableHead>
          )}
          <TableHead className='h-9 px-4 text-xs'>{t('Task status')}</TableHead>
          <TableHead className='h-9 text-xs'>{t('Delivery')}</TableHead>
          <TableHead className='h-9 text-xs'>{t('Attempts')}</TableHead>
          <TableHead className='h-9 text-xs'>{t('Last result')}</TableHead>
          <TableHead className='h-9 text-xs'>{t('Callback host')}</TableHead>
          <TableHead className='h-9 text-xs'>{t('Updated')}</TableHead>
          {props.onRedeliver && (
            <TableHead className='h-9 pr-4 text-right text-xs'>
              {t('Actions')}
            </TableHead>
          )}
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.events.map((event) => (
          <TableRow key={event.event_id} className='hover:bg-muted/30'>
            {props.showTask && (
              <TableCell className='px-4 py-2 font-mono text-xs'>
                <div>{event.task_id}</div>
                <div className='text-muted-foreground'>
                  {t('User')} #{event.user_id}
                </div>
              </TableCell>
            )}
            <TableCell className='px-4 py-2 text-xs'>
              <div className='font-medium'>{event.status}</div>
              <div
                className='text-muted-foreground font-mono text-[11px]'
                title={event.event_id}
              >
                {event.event_id}
              </div>
            </TableCell>
            <TableCell className='py-2'>
              <Badge variant={STATE_VARIANT[event.state] ?? 'outline'}>
                {t(taskCallbackStateLabel(event.state))}
              </Badge>
            </TableCell>
            <TableCell className='py-2 text-xs tabular-nums'>
              {event.attempts}/{event.max_attempts}
            </TableCell>
            <TableCell
              className='max-w-[240px] truncate py-2 text-xs'
              title={event.last_error || undefined}
            >
              {lastResult(event)}
              {event.last_latency_ms > 0 && (
                <span className='text-muted-foreground'>
                  {' '}
                  · {event.last_latency_ms} ms
                </span>
              )}
            </TableCell>
            <TableCell className='text-muted-foreground py-2 font-mono text-xs'>
              {event.callback_host || '-'}
            </TableCell>
            <TableCell
              className='text-muted-foreground py-2 text-xs whitespace-nowrap'
              title={formatTimestampToDate(event.updated_at)}
            >
              {formatTimestampRelative(event.updated_at, 'seconds', locale)}
            </TableCell>
            {props.onRedeliver && (
              <TableCell className='py-2 pr-4 text-right'>
                <Button
                  size='sm'
                  variant='outline'
                  disabled={
                    !event.redeliverable ||
                    props.redeliveringId === event.event_id
                  }
                  onClick={() => props.onRedeliver?.(event)}
                >
                  {t('Redeliver')}
                </Button>
              </TableCell>
            )}
          </TableRow>
        ))}
      </TableBody>
    </StaticDataTable>
  )
}
