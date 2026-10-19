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
import { getChannelConfigurationState } from '../channel-configuration'
import {
  CHANNEL_FORM_DEFAULT_VALUES,
  channelFormSchema,
  transformChannelToFormDefaults,
  transformFormDataToCreatePayload,
  type ChannelFormValues,
} from '../channel-form'
import {
  DEFAULT_UPSTREAM_QUOTA_PER_UNIT,
  getUpstreamQuotaPerUnit,
} from '../upstream-quota-per-unit'

function formOf(
  type: number,
  quotaPerUnit: number | undefined
): ChannelFormValues {
  return {
    ...CHANNEL_FORM_DEFAULT_VALUES,
    name: 'seedance',
    type,
    key: 'k',
    models: 'doubao-seedance-2-0-260128',
    upstream_quota_per_unit: quotaPerUnit,
  }
}

function settingsOf(type: number, quotaPerUnit: number | undefined) {
  const payload = transformFormDataToCreatePayload(formOf(type, quotaPerUnit))
  return JSON.parse(String(payload.channel.settings)) as Record<string, unknown>
}

describe('upstream quota per unit channel setting', () => {
  test('defaults to New API 500000 and stores only a custom value for Doubao video', () => {
    expect(DEFAULT_UPSTREAM_QUOTA_PER_UNIT).toBe(500000)
    expect(CHANNEL_FORM_DEFAULT_VALUES.upstream_quota_per_unit).toBe(500000)
    expect(settingsOf(CHANNEL_TYPE_DOUBAO_VIDEO, 500000)).not.toHaveProperty(
      'upstream_quota_per_unit'
    )
    expect(settingsOf(CHANNEL_TYPE_DOUBAO_VIDEO, undefined)).not.toHaveProperty(
      'upstream_quota_per_unit'
    )
    expect(
      settingsOf(CHANNEL_TYPE_DOUBAO_VIDEO, 1000000).upstream_quota_per_unit
    ).toBe(1000000)
    expect(settingsOf(CHANNEL_TYPE_OLLAMA, 1000000)).not.toHaveProperty(
      'upstream_quota_per_unit'
    )
  })

  test('round-trips from stored channel settings', () => {
    const load = (settings: string) =>
      transformChannelToFormDefaults({
        id: 1,
        name: 'seedance',
        type: CHANNEL_TYPE_DOUBAO_VIDEO,
        settings,
        channel_info: { multi_key_mode: 'random' },
      } as Parameters<typeof transformChannelToFormDefaults>[0])
    expect(
      load(JSON.stringify({ upstream_quota_per_unit: 1000000 }))
        .upstream_quota_per_unit
    ).toBe(1000000)
    expect(load('{}').upstream_quota_per_unit).toBe(500000)
  })

  test('rejects a non-positive value only on Doubao video channels', () => {
    const issues = (values: ChannelFormValues) => {
      const result = channelFormSchema.safeParse(values)
      return result.success
        ? []
        : result.error.issues.filter(
            (issue) => issue.path[0] === 'upstream_quota_per_unit'
          )
    }
    expect(issues(formOf(CHANNEL_TYPE_DOUBAO_VIDEO, 0))).toHaveLength(1)
    expect(issues(formOf(CHANNEL_TYPE_DOUBAO_VIDEO, -5))).toHaveLength(1)
    expect(issues(formOf(CHANNEL_TYPE_DOUBAO_VIDEO, 250000))).toHaveLength(0)
    expect(issues(formOf(CHANNEL_TYPE_DOUBAO_VIDEO, undefined))).toHaveLength(0)
    expect(issues(formOf(CHANNEL_TYPE_OLLAMA, 0))).toHaveLength(0)
  })

  test('marks the extra settings block as configured for a custom value', () => {
    const extraSettings = (values: ChannelFormValues) =>
      getChannelConfigurationState(values, {}, true).blocks.extraSettings
    expect(extraSettings(formOf(CHANNEL_TYPE_DOUBAO_VIDEO, 500000))).toBe(
      'idle'
    )
    expect(extraSettings(formOf(CHANNEL_TYPE_DOUBAO_VIDEO, 1000000))).toBe(
      'configured'
    )
  })

  test('reads the value used by the balance tooltip', () => {
    expect(getUpstreamQuotaPerUnit('{"upstream_quota_per_unit":250000}')).toBe(
      250000
    )
    expect(getUpstreamQuotaPerUnit('{"upstream_quota_per_unit":0}')).toBe(
      500000
    )
    expect(getUpstreamQuotaPerUnit('not json')).toBe(500000)
    expect(getUpstreamQuotaPerUnit(undefined)).toBe(500000)
    expect(getUpstreamQuotaPerUnit('')).toBe(500000)
  })
})
