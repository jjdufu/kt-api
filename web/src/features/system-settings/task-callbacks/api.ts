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
import { api } from '@/lib/api'
import { requireServerSuccess } from '@/lib/server-error-message'

export type TaskCallbackEvent = {
  event_id: string
  task_id: string
  user_id: number
  status: string
  seq: number
  state: string
  attempts: number
  max_attempts: number
  next_attempt_at: number
  last_http_status: number
  last_error: string
  last_latency_ms: number
  delivered_at: number
  created_at: number
  updated_at: number
  payload_sha256: string
  callback_host?: string
  payload?: string
  redeliverable: boolean
}

export type TaskCallbackEventFilters = {
  task_id?: string
  user_id?: number
  state?: string
  status?: string
}

type ApiEnvelope<T> = {
  success: boolean
  message?: string
  data: T
}

export type TaskCallbackEventPage = {
  page: number
  page_size: number
  total: number
  items: TaskCallbackEvent[] | null
}

export async function listTaskCallbackEvents(
  page: number,
  pageSize: number,
  filters: TaskCallbackEventFilters
): Promise<TaskCallbackEventPage> {
  const res = await api.get<ApiEnvelope<TaskCallbackEventPage>>(
    '/api/admin/task-callbacks',
    { params: { p: page, page_size: pageSize, ...filters } }
  )
  return requireServerSuccess(res.data).data
}

export async function redeliverTaskCallbackEvent(eventId: string) {
  const res = await api.post<ApiEnvelope<{ event_id: string }>>(
    `/api/admin/task-callbacks/${encodeURIComponent(eventId)}/redeliver`
  )
  return requireServerSuccess(res.data).data
}
