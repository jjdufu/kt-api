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
// Field metadata for the task callback settings form. Labels and
// descriptions are i18n source keys rendered with t().

export const TASK_CALLBACK_INJECT_MODES = [
  'only_when_customer_callback',
  'always',
  'never',
] as const

export type TaskCallbackInjectMode = (typeof TASK_CALLBACK_INJECT_MODES)[number]

export const TASK_CALLBACK_INJECT_MODE_LABEL: Record<
  TaskCallbackInjectMode,
  string
> = {
  only_when_customer_callback: 'Only when the customer set callback_url',
  always: 'For every task on enabled channels',
  never: 'Never (polling only)',
}

export const TASK_CALLBACK_BOOLEAN_FIELDS = [
  {
    name: 'enabled',
    label: 'Enable task callbacks',
    description:
      'Accept callback_url on Seedance video tasks and deliver status callbacks. When off, callback_url is ignored.',
  },
  {
    name: 'extra_headers_enabled',
    label: 'Send X-KtApi headers',
    description:
      'Add X-KtApi-Event-Id, X-KtApi-Delivery-Attempt and X-KtApi-Task-Id so receivers can deduplicate.',
  },
  {
    name: 'https_only',
    label: 'Require HTTPS callback URLs',
    description: 'Reject http:// callback URLs at submit time.',
  },
  {
    name: 'allow_private_ip',
    label: 'Allow private callback targets',
    description:
      'For testing only. Lets callbacks reach private, loopback and link-local addresses.',
  },
  {
    name: 'domain_allowlist_mode',
    label: 'Callback domain allowlist mode',
    description: 'Only accept callback URLs whose host is in the list below.',
  },
  {
    name: 'reject_unresolvable',
    label: 'Reject unresolvable callback hosts',
    description:
      'Reject a callback URL whose host does not resolve when the task is submitted.',
  },
  {
    name: 'extended_retry',
    label: 'Extended retries',
    description:
      'After the documented retries, keep retrying succeeded/failed callbacks with back-off for up to 24 hours.',
  },
  {
    name: 'signing_enabled',
    label: 'Sign callbacks',
    description:
      'Add X-KtApi-Timestamp and X-KtApi-Signature (HMAC-SHA256). Requires a signing master secret.',
  },
] as const

export const TASK_CALLBACK_NUMBER_FIELDS = [
  {
    name: 'timeout_seconds',
    label: 'Delivery timeout (seconds)',
    description: 'Per-attempt timeout. The supplier documents 5 seconds.',
    min: 1,
    max: 30,
  },
  {
    name: 'retry_interval_seconds',
    label: 'Retry interval (seconds)',
    description: 'Wait between the documented retries.',
    min: 1,
    max: 60,
  },
  {
    name: 'terminal_max_attempts',
    label: 'Attempts for succeeded/failed',
    description: 'First attempt plus retries. Other statuses are sent once.',
    min: 1,
    max: 10,
  },
  {
    name: 'workers',
    label: 'Delivery workers',
    description: 'Concurrent deliveries per instance.',
    min: 1,
    max: 128,
  },
  {
    name: 'per_host_concurrency',
    label: 'Concurrency per callback host',
    description: 'Concurrent deliveries to one receiving host.',
    min: 1,
    max: 64,
  },
  {
    name: 'retention_days',
    label: 'Delivery record retention (days)',
    description: 'Finished delivery records older than this are deleted.',
    min: 1,
    max: 90,
  },
  {
    name: 'max_url_length',
    label: 'Maximum callback URL length',
    description: 'Longer callback URLs are rejected at submit time.',
    min: 64,
    max: 8192,
  },
] as const

export type TaskCallbackSettingValues = {
  'task_callback_setting.enabled': boolean
  'task_callback_setting.extra_headers_enabled': boolean
  'task_callback_setting.https_only': boolean
  'task_callback_setting.allow_private_ip': boolean
  'task_callback_setting.domain_allowlist_mode': boolean
  'task_callback_setting.reject_unresolvable': boolean
  'task_callback_setting.extended_retry': boolean
  'task_callback_setting.signing_enabled': boolean
  'task_callback_setting.domain_list': string[]
  'task_callback_setting.allowed_ports': string[]
  'task_callback_setting.max_url_length': number
  'task_callback_setting.timeout_seconds': number
  'task_callback_setting.retry_interval_seconds': number
  'task_callback_setting.terminal_max_attempts': number
  'task_callback_setting.workers': number
  'task_callback_setting.per_host_concurrency': number
  'task_callback_setting.retention_days': number
  'task_callback_setting.upstream_public_base': string
  'task_callback_setting.upstream_inject_mode': TaskCallbackInjectMode
}

// Mirrors setting/system_setting/task_callback.go defaults.
export const DEFAULT_TASK_CALLBACK_SETTINGS: TaskCallbackSettingValues = {
  'task_callback_setting.enabled': true,
  'task_callback_setting.extra_headers_enabled': true,
  'task_callback_setting.https_only': false,
  'task_callback_setting.allow_private_ip': false,
  'task_callback_setting.domain_allowlist_mode': false,
  'task_callback_setting.reject_unresolvable': true,
  'task_callback_setting.extended_retry': false,
  'task_callback_setting.signing_enabled': false,
  'task_callback_setting.domain_list': [],
  'task_callback_setting.allowed_ports': ['80', '443', '8080', '8443'],
  'task_callback_setting.max_url_length': 2048,
  'task_callback_setting.timeout_seconds': 5,
  'task_callback_setting.retry_interval_seconds': 5,
  'task_callback_setting.terminal_max_attempts': 4,
  'task_callback_setting.workers': 16,
  'task_callback_setting.per_host_concurrency': 4,
  'task_callback_setting.retention_days': 7,
  'task_callback_setting.upstream_public_base': '',
  'task_callback_setting.upstream_inject_mode': 'only_when_customer_callback',
}

export const TASK_CALLBACK_DELIVERY_STATES = [
  'pending',
  'delivering',
  'succeeded',
  'failed',
  'superseded',
  'cancelled',
] as const

export type TaskCallbackDeliveryState =
  (typeof TASK_CALLBACK_DELIVERY_STATES)[number]

export const TASK_CALLBACK_DELIVERY_STATE_LABEL: Record<
  TaskCallbackDeliveryState,
  string
> = {
  pending: 'Pending',
  delivering: 'Delivering',
  succeeded: 'Delivered',
  failed: 'Failed',
  superseded: 'Superseded',
  cancelled: 'Cancelled',
}

export function taskCallbackStateLabel(state: string): string {
  return (
    TASK_CALLBACK_DELIVERY_STATE_LABEL[state as TaskCallbackDeliveryState] ??
    state
  )
}
