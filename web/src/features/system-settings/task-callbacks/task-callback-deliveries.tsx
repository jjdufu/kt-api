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
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { ConfirmDialog } from '@/components/confirm-dialog'
import {
  DataTablePagination,
  DataTableToolbar,
  useDataTable,
} from '@/components/data-table'
import { EmptyState } from '@/components/empty-state'
import { ErrorState } from '@/components/error-state'
import { Input } from '@/components/ui/input'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'
import { Skeleton } from '@/components/ui/skeleton'
import { handleServerError } from '@/lib/handle-server-error'

import {
  listTaskCallbackEvents,
  redeliverTaskCallbackEvent,
  type TaskCallbackEvent,
} from './api'
import {
  TASK_CALLBACK_DELIVERY_STATES,
  taskCallbackStateLabel,
} from './constants'
import { TaskCallbackEventsTable } from './task-callback-events-table'

const EMPTY_EVENTS: TaskCallbackEvent[] = []
const DELIVERIES_QUERY_KEY = ['system-settings', 'task-callbacks', 'events']

/** Admin list of task callback deliveries with filters and redelivery. */
export function TaskCallbackDeliveries() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [filters, setFilters] = useState({ state: '', task_id: '' })
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 20 })
  const [pending, setPending] = useState<TaskCallbackEvent | null>(null)

  const eventsQuery = useQuery({
    queryKey: [...DELIVERIES_QUERY_KEY, filters, pagination],
    queryFn: () =>
      listTaskCallbackEvents(pagination.pageIndex + 1, pagination.pageSize, {
        state: filters.state || undefined,
        task_id: filters.task_id.trim() || undefined,
      }),
    retry: false,
  })
  const events = eventsQuery.data?.items ?? EMPTY_EVENTS
  const { table } = useDataTable({
    data: events,
    columns: [],
    totalCount: eventsQuery.data?.total ?? 0,
    manualPagination: true,
    columnFilters: [],
    pagination,
    onPaginationChange: setPagination,
    columnVisibilityStorageKey: false,
    columnSizingStorageKey: false,
  })

  const redeliverMutation = useMutation({
    mutationFn: (event: TaskCallbackEvent) =>
      redeliverTaskCallbackEvent(event.event_id),
    onSuccess: async () => {
      toast.success(t('Callback queued for redelivery'))
      setPending(null)
      await queryClient.invalidateQueries({ queryKey: DELIVERIES_QUERY_KEY })
    },
    onError: (error) => handleServerError(error, t('Redelivery failed')),
  })

  const resetPage = () =>
    setPagination((previous) => ({ ...previous, pageIndex: 0 }))

  return (
    <div className='space-y-3'>
      <div>
        <h4 className='text-sm font-medium'>{t('Callback deliveries')}</h4>
        <p className='text-muted-foreground mt-0.5 text-xs'>
          {t(
            'Recent customer callbacks across all tasks. Failed or cancelled callbacks can be redelivered with the same event ID.'
          )}
        </p>
      </div>
      <DataTableToolbar
        table={table}
        hideViewOptions
        hasAdditionalFilters={!!filters.state || !!filters.task_id}
        onReset={() => {
          setFilters({ state: '', task_id: '' })
          resetPage()
        }}
        customSearch={
          <>
            <Input
              aria-label={t('Task ID')}
              placeholder={t('Task ID')}
              className='h-8 w-56'
              value={filters.task_id}
              onChange={(event) => {
                setFilters((previous) => ({
                  ...previous,
                  task_id: event.target.value,
                }))
                resetPage()
              }}
            />
            <NativeSelect
              aria-label={t('Delivery')}
              className='max-w-full'
              value={filters.state}
              onChange={(event) => {
                setFilters((previous) => ({
                  ...previous,
                  state: event.target.value,
                }))
                resetPage()
              }}
            >
              <NativeSelectOption value=''>
                {t('All Status')}
              </NativeSelectOption>
              {TASK_CALLBACK_DELIVERY_STATES.map((state) => (
                <NativeSelectOption key={state} value={state}>
                  {t(taskCallbackStateLabel(state))}
                </NativeSelectOption>
              ))}
            </NativeSelect>
          </>
        }
      />
      <div aria-busy={eventsQuery.isFetching}>
        {eventsQuery.isLoading && <Skeleton className='h-24 w-full' />}
        {eventsQuery.isError && (
          <ErrorState
            title={t('We could not load task callbacks.')}
            description={eventsQuery.error.message}
            onRetry={() => void eventsQuery.refetch()}
          />
        )}
        {!eventsQuery.isLoading && !eventsQuery.isError && (
          <>
            {events.length > 0 ? (
              <TaskCallbackEventsTable
                events={events}
                showTask
                redeliveringId={
                  redeliverMutation.isPending
                    ? redeliverMutation.variables?.event_id
                    : undefined
                }
                onRedeliver={setPending}
              />
            ) : (
              <EmptyState
                title={t('No task callbacks yet.')}
                className='min-h-32'
                bordered
              />
            )}
            <div className='mt-3'>
              <DataTablePagination table={table} compact />
            </div>
          </>
        )}
      </div>
      <ConfirmDialog
        open={pending !== null}
        onOpenChange={(open) => {
          if (!open && !redeliverMutation.isPending) setPending(null)
        }}
        title={t('Redeliver callback')}
        desc={t(
          'Send this callback again with a fresh attempt budget? The receiver gets the same event ID and body.'
        )}
        isLoading={redeliverMutation.isPending}
        confirmText={t('Redeliver')}
        handleConfirm={() => {
          if (pending) redeliverMutation.mutate(pending)
        }}
      />
    </div>
  )
}
