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
import { describe, expect, test } from 'vitest'

import {
  DEFAULT_TASK_CALLBACK_SETTINGS,
  TASK_CALLBACK_BOOLEAN_FIELDS,
  TASK_CALLBACK_NUMBER_FIELDS,
  taskCallbackStateLabel,
} from '../constants'

describe('task callback settings metadata', () => {
  test('defaults mirror the backend defaults', () => {
    expect(DEFAULT_TASK_CALLBACK_SETTINGS).toMatchObject({
      'task_callback_setting.enabled': true,
      'task_callback_setting.extra_headers_enabled': true,
      'task_callback_setting.https_only': false,
      'task_callback_setting.timeout_seconds': 5,
      'task_callback_setting.retry_interval_seconds': 5,
      'task_callback_setting.terminal_max_attempts': 4,
      'task_callback_setting.extended_retry': false,
      'task_callback_setting.upstream_public_base': '',
      'task_callback_setting.upstream_inject_mode':
        'only_when_customer_callback',
    })
  })

  test('every form field has a stored option', () => {
    for (const field of [
      ...TASK_CALLBACK_BOOLEAN_FIELDS,
      ...TASK_CALLBACK_NUMBER_FIELDS,
    ]) {
      expect(DEFAULT_TASK_CALLBACK_SETTINGS).toHaveProperty([
        `task_callback_setting.${field.name}`,
      ])
    }
  })

  test('delivery state labels fall back to the raw state', () => {
    expect(taskCallbackStateLabel('succeeded')).toBe('Delivered')
    expect(taskCallbackStateLabel('mystery')).toBe('mystery')
  })
})
