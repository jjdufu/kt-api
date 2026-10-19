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
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { redeliverTaskCallbackEvent } from '@/features/system-settings/task-callbacks/api'
import { TaskCallbackEventsTable } from '@/features/system-settings/task-callbacks/task-callback-events-table'
import { handleServerError } from '@/lib/handle-server-error'

import { getTaskCallbacks } from '../api'

type TaskCallbacksSectionProps = {
  taskId: string
  isAdmin: boolean
}

/**
 * Customer callback deliveries of one task (Seedance callback_url). Renders
 * nothing when the task has no callback subscription.
 */
export function TaskCallbacksSection(props: TaskCallbacksSectionProps) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const queryKey = ['usage-logs', 'task-callbacks', props.taskId]
  const callbacksQuery = useQuery({
    queryKey,
    queryFn: () => getTaskCallbacks(props.taskId),
    enabled: props.taskId !== '',
    retry: false,
    meta: { errorToast: false },
  })
  const redeliverMutation = useMutation({
    mutationFn: (eventId: string) => redeliverTaskCallbackEvent(eventId),
    onSuccess: async () => {
      toast.success(t('Callback queued for redelivery'))
      await queryClient.invalidateQueries({ queryKey })
    },
    onError: (error) => handleServerError(error, t('Redelivery failed')),
  })

  if (callbacksQuery.isLoading) {
    return <Skeleton className='h-16 w-full' />
  }
  const subscription = callbacksQuery.data?.subscription
  if (callbacksQuery.isError || !subscription) {
    return null
  }

  return (
    <section className='min-w-0 space-y-1.5'>
      <Label className='text-xs font-semibold'>{t('Callbacks')}</Label>
      <div className='bg-muted/30 min-w-0 space-y-2 rounded-md border p-2.5'>
        <div className='text-muted-foreground grid gap-1 text-xs sm:grid-cols-3'>
          <span>
            {t('Callback host')}:{' '}
            <span className='text-foreground font-mono'>
              {subscription.callback_host || '-'}
            </span>
          </span>
          <span>
            {t('URL hash')}:{' '}
            <span className='text-foreground font-mono'>
              {subscription.callback_url_hash || '-'}
            </span>
          </span>
          <span>
            {t('Upstream fast path')}:{' '}
            <span className='text-foreground'>
              {subscription.upstream_fast_path ? t('Enabled') : t('Disabled')}
            </span>
          </span>
        </div>
        {callbacksQuery.data && callbacksQuery.data.events.length > 0 ? (
          <TaskCallbackEventsTable
            events={callbacksQuery.data.events}
            redeliveringId={
              redeliverMutation.isPending
                ? redeliverMutation.variables
                : undefined
            }
            onRedeliver={
              props.isAdmin
                ? (event) => redeliverMutation.mutate(event.event_id)
                : undefined
            }
          />
        ) : (
          <p className='text-muted-foreground text-xs'>
            {t('No callbacks have been sent for this task yet.')}
          </p>
        )}
      </div>
    </section>
  )
}
