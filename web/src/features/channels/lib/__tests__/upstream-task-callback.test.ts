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

import { CHANNEL_TYPE_DOUBAO_VIDEO, CHANNEL_TYPE_OLLAMA } from '../../constants'
import {
  CHANNEL_FORM_DEFAULT_VALUES,
  transformChannelToFormDefaults,
  transformFormDataToCreatePayload,
} from '../channel-form'

function settingsOf(type: number, enabled: boolean) {
  const payload = transformFormDataToCreatePayload({
    ...CHANNEL_FORM_DEFAULT_VALUES,
    name: 'seedance',
    type,
    key: 'k',
    models: 'doubao-seedance-2-0-260128',
    upstream_task_callback_enabled: enabled,
  })
  return JSON.parse(String(payload.channel.settings)) as Record<string, unknown>
}

describe('upstream task callback channel setting', () => {
  test('is off by default and saved only for Doubao video channels', () => {
    expect(CHANNEL_FORM_DEFAULT_VALUES.upstream_task_callback_enabled).toBe(
      false
    )
    expect(
      settingsOf(CHANNEL_TYPE_DOUBAO_VIDEO, true).upstream_task_callback_enabled
    ).toBe(true)
    expect(settingsOf(CHANNEL_TYPE_DOUBAO_VIDEO, false)).not.toHaveProperty(
      'upstream_task_callback_enabled'
    )
    expect(settingsOf(CHANNEL_TYPE_OLLAMA, true)).not.toHaveProperty(
      'upstream_task_callback_enabled'
    )
  })

  test('round-trips from stored channel settings', () => {
    const defaults = transformChannelToFormDefaults({
      id: 1,
      name: 'seedance',
      type: CHANNEL_TYPE_DOUBAO_VIDEO,
      settings: JSON.stringify({ upstream_task_callback_enabled: true }),
      channel_info: { multi_key_mode: 'random' },
    } as Parameters<typeof transformChannelToFormDefaults>[0])
    expect(defaults.upstream_task_callback_enabled).toBe(true)
  })
})
